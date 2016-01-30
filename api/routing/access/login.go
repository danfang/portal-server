package access

import (
	"encoding/hex"
	"portal-server/api/errs"
	"portal-server/api/routing"
	"portal-server/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// A PasswordLogin is a JSON structure for user logins via email and password.
//
// swagger:parameters login
type PasswordLogin struct {
	// in: body
	// required: true
	Body passwordLogin `json:"password_login"`
}

type passwordLogin struct {
	// unique: true
	// required: true
	Email string `json:"email" valid:"required,email"`

	// minimum length: 6
	// maximum length: 50
	// required: true
	Password string `json:"password" valid:"required,length(6|50)"`
}

// LoginEndpoint handles a POST request for a user to login via email and password.
func (r Router) LoginEndpoint(c *gin.Context) {
	var body passwordLogin
	if !routing.ValidateJSON(c, &body) {
		return
	}

	var user model.User
	if r.Db.Where("email = ?", body.Email).First(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, routing.RenderError(errs.ErrInvalidLogin))
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, routing.RenderError(errs.ErrInvalidLogin))
		return
	}

	split := strings.Split(user.Password, ":")
	password := split[0]
	salt, err := hex.DecodeString(split[1])

	// Get the salt value
	if err != nil {
		routing.InternalServiceError(c, err)
		return
	}

	// Compare the two password hashes
	if hashPassword(body.Password, salt) != password {
		c.JSON(http.StatusBadRequest, routing.RenderError(errs.ErrInvalidLogin))
		return
	}

	userToken, err := createUserToken(r.Db, &user)
	if err != nil {
		routing.InternalServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		UserUUID:  user.UUID,
		UserToken: userToken,
	})
}
