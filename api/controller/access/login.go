package access

import (
	"encoding/hex"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/model"
	"strings"

	"github.com/gin-gonic/gin"
	"portal-server/api/controller/context"
)

type passwordLogin struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6|50)"`
}

// LoginEndpoint handles a POST request for a user to login via email and password.
func LoginEndpoint(c *gin.Context) {
	var body passwordLogin
	if !controller.ValidJSON(c, &body) {
		return
	}

	store := context.StoreFromContext(c)
	user, found := store.Users().FindUser(&model.User{Email: body.Email})
	if !found || user.Password == "" {
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

	userToken, err := createUserToken(store, user)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		UserUUID:  user.UUID,
		UserToken: userToken,
	})
}
