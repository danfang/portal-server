package routing

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const validatorResponse = "passed"

type testJSON struct {
	A string `valid:"required,email"`
	B bool   `valid:"required"`
	C int    `valid:"required"`
}

var validator *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	validator = gin.New()
	validator.GET("/", func(c *gin.Context) {
		var json testJSON
		if !ValidateJSON(c, &json) {
			return
		}
		c.String(http.StatusOK, validatorResponse)
	})
}

func TestValidateJSON_BadFormat(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(`{"broken_json": }`))
	w := httptest.NewRecorder()
	validator.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	expected, _ := json.Marshal(gin.H{
		"error":  "invalid_json",
		"reason": "invalid character '}' looking for beginning of value",
	})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestValidateJSON_MissingRequired(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(`{"B": true, "C": 42}`))
	w := httptest.NewRecorder()
	validator.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	expected, _ := json.Marshal(gin.H{
		"error":  "invalid_json",
		"reason": "A: non zero value required;",
	})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestValidateJSON_IncorrectType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(`{"A": true, "B": true, "C": 42}`))
	w := httptest.NewRecorder()
	validator.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	expected, _ := json.Marshal(gin.H{
		"error":  "invalid_json",
		"reason": "json: cannot unmarshal bool into Go value of type string",
	})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestValidateJSON_FailValidation(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(`{"A": "not_an_email", "B": true, "C": 42}`))
	w := httptest.NewRecorder()
	validator.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	expected, _ := json.Marshal(gin.H{
		"error":  "invalid_json",
		"reason": "A: not_an_email does not validate as email;",
	})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestValidateJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(`{"A": "email@email.com", "B": true, "C": 42}`))
	w := httptest.NewRecorder()
	validator.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.Equal(t, validatorResponse, w.Body.String())
}

func TestInternalServerError(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		InternalServiceError(c, errors.New("we_broke"))
	})
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	response, _ := json.Marshal(map[string]string{
		"error": "internal_server_error",
	})
	assert.JSONEq(t, string(response), w.Body.String())
}
