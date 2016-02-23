package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/model"
	"portal-server/store"

	"github.com/gin-gonic/gin"
)

type addContacts struct {
	Contacts []contact `json:"contacts"`
}

type contact struct {
	UUID         string         `json:"cid" valid:"required,uuid"`
	Name         string         `json:"name" valid:"required"`
	PhoneNumbers []contactPhone `json:"phone_numbers" valid:"required"`
}

type contactPhone struct {
	Number string `json:"number" valid:"required"`
	Name   string `json:"name" valid:"required"`
}

// AddContactsEndpoint allows users to upload their contacts
func AddContactsEndpoint(c *gin.Context) {
	var body addContacts
	if !controller.ValidJSON(c, &body) {
		return
	}
	user := context.UserFromContext(c)
	s := context.StoreFromContext(c)
	s.Transaction(func(store store.Store) error {
		for _, c := range body.Contacts {
			phoneNumbers := make([]model.ContactPhone, 0, len(c.PhoneNumbers))
			for _, p := range c.PhoneNumbers {
				number := model.ContactPhone{
					Number: p.Number,
					Name:   p.Name,
				}
				phoneNumbers = append(phoneNumbers, number)
			}
			if err := store.Contacts().CreateContact(&model.Contact{
				UserID:       user.ID,
				UUID:         c.UUID,
				Name:         c.Name,
				PhoneNumbers: phoneNumbers,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	c.JSON(http.StatusOK, controller.RenderSuccess(true))
}
