package access

import (
	"encoding/hex"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type passwordLogin struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6|50)"`
}

// LoginEndpoint handles a POST request for a user to login via email and password.
func (r Router) LoginEndpoint(c *gin.Context) {
	var body passwordLogin
	if !controller.ValidJSON(c, &body) {
		return
	}

	var user model.User
	if r.Db.Where("email = ?", body.Email).First(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, controller.RenderError(errs.ErrInvalidLogin))
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, controller.RenderError(errs.ErrInvalidLogin))
		return
	}

	split := strings.Split(user.Password, ":")
	password := split[0]
	salt, err := hex.DecodeString(split[1])

	// Get the salt value
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}

	// Compare the two password hashes
	if hashPassword(body.Password, salt) != password {
		c.JSON(http.StatusBadRequest, controller.RenderError(errs.ErrInvalidLogin))
		return
	}

	userToken, err := createUserToken(r.Db, &user)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		UserUUID:  user.UUID,
		UserToken: userToken,
	})
}
