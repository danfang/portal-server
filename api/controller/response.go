package controller

import (
	"net/http"
	"portal-server/api/errs"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error"`
}

type DetailError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func InternalServiceError(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(http.StatusInternalServerError, RenderError(errs.ErrInternal))
}

func RenderError(e error) interface{} {
	return Error{e.Error()}
}

func RenderSuccess() interface{} {
	return SuccessResponse{true}
}
