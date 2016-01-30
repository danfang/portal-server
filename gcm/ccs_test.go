package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoogleCCS_BadCredentials(t *testing.T) {
	ccs := GoogleCCS{"id", "key"}
	err := ccs.Listen(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not-authorized")
}

func TestGoogleCCS(t *testing.T) {
	ccs := GoogleCCS{senderID, apiKey}
	stop := make(chan bool)
	go func() {
		if ccs.Listen(nil, stop) != nil {
			t.Fail()
		}
	}()
	time.Sleep(1 * time.Second)
	stop <- true
}
