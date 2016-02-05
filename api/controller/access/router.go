package access

import (
	"net/http"

	"portal-server/store"
)

// Router handles access routes for user creation and authentication.
type Router struct {
	Store      store.Store
	HTTPClient *http.Client
}
