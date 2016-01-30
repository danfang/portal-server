package util

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// TestHTTP allows a test to mock an HTTP response for any WebClient. This is based on
// http://keighl.com/post/mocking-http-responses-in-golang/
func TestHTTP(requestTest func(*http.Request), responseCode int, output string) (*httptest.Server, *WebClient) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTest(r)
		w.WriteHeader(responseCode)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, output)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &WebClient{server.URL, httpClient}

	return server, client
}
