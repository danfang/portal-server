package routing

import (
	"portal-server/api/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// An ErrorResponse wraps a bad request with an error code string.
//
// swagger:response error
type ErrorResponse struct {
	// in: body
	Body Error
}

// An Error contains an error JSON object containing an error.
type Error struct {
	Error string `json:"error"`
}

// A DetailErrorResponse gives both the error and reason for the error.
//
// swagger:response detailError
type DetailErrorResponse struct {
	// in: body
	Body DetailError
}

// A DetailError contains a JSON object with both an error and an error reason.
type DetailError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// A SuccessResponse denotes whether or not an action was successful.
//
// swagger:response success
type SuccessResponse struct {
	// in: body
	Body successResponse
}

type successResponse struct {
	Success bool `json:"success"`
}

// InternalServiceError records and writes an internal service error
// response.
func InternalServiceError(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(http.StatusInternalServerError, RenderError(errs.ErrInternal))
}

// RenderError generates JSON output for an error.
func RenderError(e error) interface{} {
	return Error{e.Error()}
}

// RenderSuccess generates JSON output for a successful operation.
func RenderSuccess() interface{} {
	return successResponse{true}
}
