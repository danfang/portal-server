package access

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Correct(t *testing.T) {
	salt := []byte{255, 10, 25, 16}
	assert.Equal(t, hashPassword("password", salt), hashPassword("password", salt))
}

func TestHashPassword_BadPassword(t *testing.T) {
	salt := []byte{255, 10, 25, 16}
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("PasSworD", salt))
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("password1", salt))
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("", salt))
}

func TestHashPassword_BadSalt(t *testing.T) {
	salt := []byte{255, 10, 25, 16}
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("password", []byte{10, 10, 10}))
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("password", []byte{}))
	assert.NotEqual(t, hashPassword("password", salt), hashPassword("password", nil))
}
