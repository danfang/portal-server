package controller

import (
	"encoding/json"
	"portal-server/api/errs"
	"testing"

	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("RenderError", func() {
		testCases := []error{
			errs.ErrInvalidLogin,
			errs.ErrInternal,
			errs.ErrInvalidJSON,
			errs.ErrMissingHeaders,
			errs.ErrInvalidUserToken,
			errs.ErrDuplicateEmail,
			errs.ErrUnsupportedAccountType,
			errs.ErrInvalidRegistrationToken,
			errs.GCMError("gcm_error"),
		}
		g.It("Should render all errors correctly", func() {
			for _, err := range testCases {
				expected, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
				actual, _ := json.Marshal(RenderError(err))
				assert.JSONEq(t, string(expected), string(actual))
			}
		})
	})

	g.Describe("RenderSuccess", func() {
		g.It("Should render success = true correctly", func() {
			actual, _ := json.Marshal(RenderSuccess(true))
			assert.JSONEq(t, "{ \"success\": true }", string(actual))
		})

		g.It("Should render success = false correctly", func() {
			actual, _ := json.Marshal(RenderSuccess(false))
			assert.JSONEq(t, "{ \"success\": false }", string(actual))
		})
	})

}
