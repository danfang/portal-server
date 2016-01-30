package routing

import (
	"encoding/json"
	"github.com/danfang/portal-server/api/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderError(t *testing.T) {
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

	for _, err := range testCases {
		expected, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
		actual, _ := json.Marshal(RenderError(err))
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestRenderSuccess(t *testing.T) {
	expected, _ := json.Marshal(map[string]interface{}{
		"success": true,
	})
	actual, _ := json.Marshal(RenderSuccess())
	assert.JSONEq(t, string(expected), string(actual))
}
