package access

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type VerificationToken struct {
	Token string `json:"token"`
}

// VerifyUserEndpoint handles a GET request that consumes a user's verification token
// for users who registered with an email and password.
func (r Router) VerifyUserEndpoint(c *gin.Context) {
	user, err := checkVerificationToken(r.Db, c.Param("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, controller.RenderError(err))
		return
	}
	user.Verified = true
	if err := r.Db.Save(&user).Error; err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, controller.RenderSuccess())
}

func checkVerificationToken(db *gorm.DB, param string) (*model.User, error) {
	var token model.VerificationToken

	// Check for existing token
	if db.Where(model.VerificationToken{Token: param}).First(&token).RecordNotFound() {
		return nil, errs.ErrInvalidVerificationToken
	}

	defer db.Delete(&token)

	// Expired token
	if time.Now().After(token.ExpiresAt) {
		return nil, errs.ErrExpiredVerificationToken
	}

	// Check for existing user account
	var user model.User
	if err := db.Model(&token).Related(&user).Error; err != nil {
		return nil, errs.ErrInvalidVerificationToken
	}

	return &user, nil
}
