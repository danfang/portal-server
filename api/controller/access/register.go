package access

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/api/errs"
	"portal-server/model"
	"portal-server/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type passwordRegistration struct {
	Email       string `json:"email" valid:"required,email"`
	Password    string `json:"password" valid:"required,length(6|50)"`
	FirstName   string `json:"first_name" valid:"length(1|20)"`
	LastName    string `json:"last_name" valid:"length(1|20)"`
	PhoneNumber string `json:"phone_number" valid:"matches(^\+[0-9]{10,12}$)"`
}

// RegisterEndpoint handles a POST request to register a new user via
// email and password.
func RegisterEndpoint(c *gin.Context) {
	var body passwordRegistration
	if !controller.ValidJSON(c, &body) {
		return
	}

	var user *model.User
	var verificationToken string
	s := context.StoreFromContext(c)
	s.Transaction(func(store store.Store) error {
		var err error

		// Check unique email
		if store.Users().UserCount(&model.User{Email: body.Email}) >= 1 {
			err = errs.ErrDuplicateEmail
			c.JSON(http.StatusBadRequest, controller.RenderError(err))
			return err
		}

		user, err = createDefaultUser(store, &body)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}

		verificationToken, err = createVerificationToken(store, user)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}
		return nil
	})
	// Send confirmation email to user
	sendTokenToUser(user.Email, verificationToken)
	c.JSON(http.StatusOK, controller.RenderSuccess())
}

func createDefaultUser(store store.Store, body *passwordRegistration) (*model.User, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	password := hashPassword(body.Password, salt)
	user := &model.User{
		UUID:      uuid.NewV4().String(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  password + ":" + hex.EncodeToString(salt),
		Verified:  false,
	}
	if err := store.Users().CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func createVerificationToken(store store.Store, user *model.User) (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	newToken := &model.VerificationToken{
		User:      *user,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
		Token:     hex.EncodeToString(token),
	}
	if err := store.VerificationTokens().CreateToken(newToken); err != nil {
		return "", err
	}
	return newToken.Token, nil
}
