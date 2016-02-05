package access

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type PasswordRegistration struct {
	Body passwordRegistration `json:"password_registration"`
}

type passwordRegistration struct {
	Email       string `json:"email" valid:"required,email"`
	Password    string `json:"password" valid:"required,length(6|50)"`
	FirstName   string `json:"first_name" valid:"length(1|20)"`
	LastName    string `json:"last_name" valid:"length(1|20)"`
	PhoneNumber string `json:"phone_number" valid:"matches(^\+[0-9]{10,12}$)"`
}

// RegisterEndpoint handles a POST request to register a new user via
// email and password.
func (r Router) RegisterEndpoint(c *gin.Context) {
	var body passwordRegistration
	if !controller.ValidJSON(c, &body) {
		return
	}
	tx := r.Db.Begin()

	var count int
	if r.Db.Model(model.User{}).Where(model.User{Email: body.Email}).Count(&count); count >= 1 {
		c.JSON(http.StatusBadRequest, controller.RenderError(errs.ErrDuplicateEmail))
		return
	}

	user, err := createDefaultUser(tx, &body)
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	token, err := createVerificationToken(tx, user)
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	tx.Commit()
	sendTokenToUser(user.Email, token)
	c.JSON(http.StatusOK, controller.RenderSuccess())
}

func createDefaultUser(db *gorm.DB, body *passwordRegistration) (*model.User, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	password := hashPassword(body.Password, salt)
	user := model.User{
		UUID:        uuid.NewV4().String(),
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    password + ":" + hex.EncodeToString(salt),
		Verified:    false,
		PhoneNumber: body.PhoneNumber,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func createVerificationToken(db *gorm.DB, user *model.User) (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	newToken := model.VerificationToken{
		User:      *user,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
		Token:     hex.EncodeToString(token),
	}
	if err := db.Create(&newToken).Error; err != nil {
		return "", err
	}
	return newToken.Token, nil
}
