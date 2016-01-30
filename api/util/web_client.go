package util

import (
	"net/http"
)

// A WebClient makes web requests based on a given
// base URL endpoint and a given http Client.
type WebClient struct {
	BaseURL    string
	HTTPClient *http.Client
}
