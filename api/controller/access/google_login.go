package access

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/api/util"
	"portal-server/model"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"portal-server/store"
)

var googleOAuthEndpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo"

type googleLogin struct {
	IDToken string `json:"id_token" valid:"required"`
}

// GoogleLoginEndpoint handles a POST request to login or register with a Google account.
func (r Router) GoogleLoginEndpoint(c *gin.Context) {
	var body googleLogin
	if !controller.ValidJSON(c, &body) {
		return
	}

	client := &util.WebClient{
		BaseURL:    googleOAuthEndpoint,
		HTTPClient: r.HTTPClient,
	}

	googleUser, err := util.GetGoogleUser(client, body.IDToken)

	// Check for errors with the Google user
	switch {
	case err == errs.ErrInvalidGoogleIDToken:
		c.JSON(http.StatusBadRequest, controller.RenderError(err))
		return
	case err == errs.ErrGoogleOAuthUnavailable:
		c.JSON(http.StatusInternalServerError, controller.RenderError(err))
		return
	case err != nil:
		controller.InternalServiceError(c, err)
		return
	}

	if googleUser.EmailVerified == "false" {
		c.JSON(http.StatusBadRequest, controller.RenderError(errs.ErrGoogleAccountNotVerified))
		return
	}

	var user *model.User
	var userToken string
	r.Store.Transaction(func(tx store.Store) error {
		user, err = createLinkedGoogleAccount(tx, googleUser)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}

		userToken, err = createUserToken(tx, user)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}
		return nil
	})
	c.JSON(http.StatusOK, loginResponse{
		UserUUID:  user.UUID,
		UserToken: userToken,
	})
}

func createLinkedGoogleAccount(store store.Store, googleUser *util.GoogleUser) (*model.User, error) {
	var user *model.User
	account, err := store.LinkedAccounts().FindAccount(&model.LinkedAccount{
		AccountID: googleUser.Sub,
		Type:      model.LinkedAccountTypeGoogle,
	})
	if err != nil {
		// Create new user account, if none exists.
		user, err = store.Users().FindOrCreateUser(&model.User{Email: googleUser.Email}, &model.User{
			UUID:      uuid.NewV4().String(),
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
			Email:     googleUser.Email,
			Verified:  true,
		})
		if err != nil {
			return nil, err
		}
		// Disable password login and verify user
		user.Password = ""
		user.Verified = true
		if err := store.Users().SaveUser(user); err != nil {
			return nil, err
		}
		// Create linked account from the user account.
		if err := store.LinkedAccounts().CreateAccount(&model.LinkedAccount{
			User:      *user,
			AccountID: googleUser.Sub,
			Type:      model.LinkedAccountTypeGoogle,
		}); err != nil {
			return nil, err
		}
		return user, nil
	}
	user, err = store.LinkedAccounts().GetRelatedUser(account)
	if err != nil {
		return nil, err
	}
	return user, nil
}
