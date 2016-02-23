package store

import (
	"portal-server/model"
	"testing"

	"github.com/franela/goblin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestContactStore(t *testing.T) {
	var db *gorm.DB
	var store contactStore
	g := goblin.Goblin(t)

	g.Describe("ContactStore", func() {
		g.BeforeEach(func() {
			db = GetTestDB()
			store = contactStore{db}
		})

		g.AfterEach(func() {
			TeardownTestDB(db)
		})

		g.It("CreateContacts", func() {
			user := model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			contact1 := model.Contact{
				User: user,
				Name: "contact1",
				UUID: uuid.NewV4().String(),
				PhoneNumbers: []model.ContactPhone{
					{
						Name:   "home",
						Number: "1",
					},
					{
						Name:   "cell",
						Number: "2",
					},
				},
			}
			contact2 := model.Contact{
				User: user,
				Name: "contact2",
				UUID: uuid.NewV4().String(),
				PhoneNumbers: []model.ContactPhone{
					{
						Name:   "home2",
						Number: "3",
					},
					{
						Name:   "cell2",
						Number: "4",
					},
				},
			}
			contact3 := model.Contact{
				User: user,
				Name: "contact3",
				UUID: uuid.NewV4().String(),
				PhoneNumbers: []model.ContactPhone{
					{
						Name:   "home3",
						Number: "5",
					},
					{
						Name:   "cell3",
						Number: "6",
					},
				},
			}

			store.CreateContact(&contact1)
			store.CreateContact(&contact2)
			store.CreateContact(&contact3)

			var contactCount int
			db.Model(&model.Contact{}).Count(&contactCount)
			assert.Equal(t, 3, contactCount)

			var phoneCount int
			db.Model(&model.ContactPhone{}).Count(&phoneCount)
			assert.Equal(t, 6, phoneCount)

			assertContact(t, db, &user, "contact1", "home", "1")
			assertContact(t, db, &user, "contact1", "cell", "2")
			assertContact(t, db, &user, "contact2", "home2", "3")
			assertContact(t, db, &user, "contact2", "cell2", "4")
			assertContact(t, db, &user, "contact3", "home3", "5")
			assertContact(t, db, &user, "contact3", "cell3", "6")
		})
	})
}

func assertContact(t *testing.T, db *gorm.DB, user *model.User, contactName, phoneName, phoneNumber string) {
	var contact model.Contact
	db.Where(&model.Contact{
		UserID: user.ID,
		Name:   contactName,
	}).First(&contact)

	assert.Equal(t, user.ID, contact.UserID)
	assert.NotNil(t, uuid.FromStringOrNil(contact.UUID))

	var phone model.ContactPhone
	db.Where(&model.ContactPhone{
		ContactID: contact.ID,
		Name:      phoneName,
	}).First(&phone)

	assert.Equal(t, phoneNumber, phone.Number)
}
