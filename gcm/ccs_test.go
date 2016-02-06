package main

import (
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoogleCCS(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Google CCS", func() {

		g.It("Should error on bad credentials", func() {
			ccs := GoogleCCS{"id", "key"}
			err := ccs.Listen(nil, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not-authorized")
		})

		g.It("Should not error on valid credentials", func() {
			ccs := GoogleCCS{senderID, apiKey}
			stop := make(chan bool)
			go func() {
				if ccs.Listen(nil, stop) != nil {
					t.Fail()
				}
			}()
			time.Sleep(1 * time.Second)
			stop <- true
		})
	})
}
