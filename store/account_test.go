package store

import (
	"portal-server/model"
	"testing"

	"github.com/franela/goblin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestLinkedAccountStore(t *testing.T) {
	var db *gorm.DB
	var store linkedAccountStore
	g := goblin.Goblin(t)

	g.Describe("LinkedAccountStore", func() {
		g.BeforeEach(func() {
			db = GetTestDB()
			store = linkedAccountStore{db}
		})

		g.AfterEach(func() {
			TeardownTestDB(db)
		})

		g.It("CreateAccount", func() {
			user := model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			store.CreateAccount(&model.LinkedAccount{
				User:      user,
				Type:      model.LinkedAccountTypeGoogle,
				AccountID: "1234",
			})
			var acc model.LinkedAccount
			db.Where(model.LinkedAccount{UserID: user.ID}).First(&acc)
			assert.NotNil(t, acc)
			assert.Equal(t, "1234", acc.AccountID)
			assert.Equal(t, user.ID, acc.UserID)
			assert.Equal(t, model.LinkedAccountTypeGoogle, acc.Type)
		})

		g.It("FindAccount", func() {
			user := model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			account := model.LinkedAccount{
				User:      user,
				Type:      model.LinkedAccountTypeGoogle,
				AccountID: "1234",
			}
			db.Create(&account)
			acc, found := store.FindAccount(&model.LinkedAccount{UserID: user.ID})
			assert.True(t, found)
			assert.Equal(t, "1234", acc.AccountID)
			assert.Equal(t, user.ID, acc.UserID)
			assert.Equal(t, model.LinkedAccountTypeGoogle, acc.Type)

			_, found = store.FindAccount(&model.LinkedAccount{UserID: user.ID + 1})
			assert.False(t, found)
		})

		g.It("GetRelatedUser", func() {
			user := model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			account := model.LinkedAccount{
				User:      user,
				Type:      model.LinkedAccountTypeGoogle,
				AccountID: "1234",
			}
			db.Create(&account)
			related, err := store.GetRelatedUser(&account)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, related.ID)
			assert.Equal(t, user.UUID, related.UUID)
			assert.Equal(t, user.Email, related.Email)
		})

		g.It("GetCount", func() {
			user := model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			count := 5
			for i := 0; i < count; i++ {
				account := model.LinkedAccount{
					User:      user,
					Type:      model.LinkedAccountTypeGoogle,
					AccountID: string(i),
				}
				db.Create(&account)
			}
			assert.Equal(t, count, store.GetCount(&model.LinkedAccount{UserID: user.ID}))
		})
	})
}
