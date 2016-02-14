package store

import (
	"errors"
	"portal-server/model"
	"testing"

	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	var store Store
	g := goblin.Goblin(t)

	g.Describe("Store transaction", func() {
		g.BeforeEach(func() {
			store = GetTestStore()
		})

		g.AfterEach(func() {
			TeardownTestStore(store)
		})

		g.It("Should rollback all changes on error", func() {
			store.Transaction(func(s Store) error {
				user := model.User{
					UUID:  "1",
					Email: "test@portal.com",
				}
				s.Users().CreateUser(&user)
				token := model.UserToken{
					User:  user,
					Token: "token",
				}
				s.UserTokens().CreateToken(&token)
				return errors.New("an_error")
			})
			_, found := store.Users().FindUser(&model.User{})
			assert.False(t, found)

			_, found = store.UserTokens().FindToken(&model.UserToken{})
			assert.False(t, found)
		})

		g.It("Should commit all changes on no error", func() {
			store.Transaction(func(s Store) error {
				user := model.User{
					UUID:  "1",
					Email: "test@portal.com",
				}
				s.Users().CreateUser(&user)
				token := model.UserToken{
					User:  user,
					Token: "token",
				}
				s.UserTokens().CreateToken(&token)
				return nil
			})
			_, found := store.Users().FindUser(&model.User{})
			assert.True(t, found)

			_, found = store.UserTokens().FindToken(&model.UserToken{})
			assert.True(t, found)
		})
	})
}
