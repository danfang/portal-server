package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/model"
	"portal-server/store"

	"github.com/gin-gonic/gin"
)

type contactsJson struct {
	Contacts []model.Contact `json:"contacts" valid:"required"`
}

// AddContactsEndpoint allows users to upload their contacts
func AddContactsEndpoint(c *gin.Context) {
	var body contactsJson
	if !controller.ValidJSON(c, &body) {
		return
	}
	user := context.UserFromContext(c)
	s := context.StoreFromContext(c)
	s.Transaction(func(store store.Store) error {
		for _, c := range body.Contacts {
			c.UserID = user.ID
			if err := store.Contacts().CreateContact(&c); err != nil {
				return err
			}
		}
		return nil
	})
	c.JSON(http.StatusOK, controller.RenderSuccess(true))
}
