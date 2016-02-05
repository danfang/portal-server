package routing

import (
	"net/http"
	"portal-server/api/errs"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// ValidateJSON writes a response if there are JSON marshalling
// or JSON validation errors. Returns true if given JSON is valid.
func ValidJSON(c *gin.Context, json interface{}) bool {
	if err := c.BindJSON(json); err != nil {
		c.JSON(http.StatusBadRequest, DetailError{
			Error:  errs.ErrInvalidJSON.Error(),
			Reason: err.Error(),
		})
		return false
	}

	if _, err := govalidator.ValidateStruct(json); err != nil {
		c.JSON(http.StatusBadRequest, DetailError{
			Error:  errs.ErrInvalidJSON.Error(),
			Reason: err.Error(),
		})
		return false
	}
	return true
}
