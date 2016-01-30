package access

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/danfang/portal-server/api/errs"
	"github.com/danfang/portal-server/api/routing"
	"github.com/danfang/portal-server/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

// A PasswordRegistration is a JSON structure for user registration.
// via email and password.
//
// swagger:parameters register
type PasswordRegistration struct {
	// in: body
	// required: true
	Body passwordRegistration `json:"password_registration"`
}

type passwordRegistration struct {
	// unique: true
	// required: true
	Email string `json:"email" valid:"required,email"`

	// minimum length: 6
	// maximum length: 50
	// required: true
	Password string `json:"password" valid:"required,length(6|50)"`

	// minimum length: 1
	// maximum length: 20
	// required: true
	FirstName string `json:"first_name" valid:"length(1|20)"`

	// minimum length: 1
	// maximum length: 20
	// required: true
	LastName string `json:"last_name" valid:"length(1|20)"`

	// pattern: ^\+[0-9]{10,12}$
	PhoneNumber string `json:"phone_number" valid:"matches(^\+[0-9]{10,12}$)"`
}

// RegisterEndpoint handles a POST request to register a new user via
// email and password.
func (r Router) RegisterEndpoint(c *gin.Context) {
	var body passwordRegistration
	if !routing.ValidateJSON(c, &body) {
		return
	}
	tx := r.Db.Begin()

	var count int
	if r.Db.Model(model.User{}).Where(model.User{Email: body.Email}).Count(&count); count >= 1 {
		c.JSON(http.StatusBadRequest, routing.RenderError(errs.ErrDuplicateEmail))
		return
	}

	user, err := createDefaultUser(tx, &body)
	if err != nil {
		tx.Rollback()
		routing.InternalServiceError(c, err)
		return
	}

	token, err := createVerificationToken(tx, user)
	if err != nil {
		tx.Rollback()
		routing.InternalServiceError(c, err)
		return
	}

	tx.Commit()
	sendTokenToUser(user.Email, token)
	c.JSON(http.StatusOK, routing.RenderSuccess())
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
