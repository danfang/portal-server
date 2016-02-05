package access

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/api/util"
	"portal-server/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
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

	tx := r.Db.Begin()
	user, err := createLinkedGoogleAccount(tx, googleUser)
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	userToken, err := createUserToken(tx, user)
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, loginResponse{
		UserUUID:  user.UUID,
		UserToken: userToken,
	})
}

func createLinkedGoogleAccount(db *gorm.DB, googleUser *util.GoogleUser) (*model.User, error) {
	var account model.LinkedAccount
	var user model.User
	// Check if account is already linked.
	if db.Where(model.LinkedAccount{
		AccountID: googleUser.Sub,
		Type:      model.LinkedAccountTypeGoogle,
	}).First(&account).RecordNotFound() {
		// Create new user account, if none exists.
		if err := db.Where(model.User{Email: googleUser.Email}).Attrs(model.User{
			UUID:      uuid.NewV4().String(),
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
			Email:     googleUser.Email,
			Verified:  true,
		}).FirstOrCreate(&user).Error; err != nil {
			return nil, err
		}
		// Disable password login and verify user
		user.Password = ""
		user.Verified = true
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
		// Create linked account from the user account.
		if err := db.Create(&model.LinkedAccount{
			User:      user,
			AccountID: googleUser.Sub,
			Type:      model.LinkedAccountTypeGoogle,
		}).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	db.Model(&account).Related(&user)
	return &user, nil
}
