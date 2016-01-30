package access

import (
	"github.com/danfang/portal-server/api/errs"
	"github.com/danfang/portal-server/api/routing"
	"github.com/danfang/portal-server/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

// A VerificationToken is a token generated on account creation and emailed
// to a given user.
//
// swagger:parameters verifyToken
type VerificationToken struct {
	// in: path
	// required: true
	Token string `json:"token"`
}

// VerifyUserEndpoint handles a GET request that consumes a user's verification token
// for users who registered with an email and password.
func (r Router) VerifyUserEndpoint(c *gin.Context) {
	user, err := checkVerificationToken(r.Db, c.Param("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, routing.RenderError(err))
		return
	}
	user.Verified = true
	if err := r.Db.Save(&user).Error; err != nil {
		routing.InternalServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, routing.RenderSuccess())
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
