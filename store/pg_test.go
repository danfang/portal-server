package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStore(t *testing.T) {
	s := GetStore("postgres", "postgres", "password")
	assert.NotNil(t, s)
}
