package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"

	"github.com/gin-gonic/gin"
)

// GetContactsEndpoint allows users to retrieve their contacts
func GetContactsEndpoint(c *gin.Context) {
	user := context.UserFromContext(c)
	s := context.StoreFromContext(c)
	contacts, err := s.Contacts().GetContactsByUser(user)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, contactsJson{Contacts: contacts})
}
