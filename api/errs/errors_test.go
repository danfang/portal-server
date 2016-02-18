package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCMError(t *testing.T) {
	var err error = GCMError("invalid registration ids")
	err, isGCMError := err.(GCMError)
	assert.True(t, isGCMError)
	assert.EqualError(t, err, "invalid registration ids")
}
