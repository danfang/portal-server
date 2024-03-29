package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	g := goblin.Goblin(t)
	api := API(store.GetTestStore(), http.DefaultClient)

	g.Describe("API routes", func() {

		g.It("Should not find anything for GET /", func() {
			req, _ := http.NewRequest("GET", "/v1", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		g.It("Should allow a POST /register", func() {
			req, _ := http.NewRequest("POST", "/v1/register", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		g.It("Should allow a POST /login", func() {
			req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		g.It("Should allow a POST /login/google", func() {
			req, _ := http.NewRequest("POST", "/v1/login/google", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		g.It("Should not allow a GET /verify/ without a token", func() {
			req, _ := http.NewRequest("GET", "/v1/verify/", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		g.It("Should allow a GET /verify/:token", func() {
			req, _ := http.NewRequest("GET", "/v1/verify/abc", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		g.It("Should allow a POST /user/devices", func() {
			req, _ := http.NewRequest("POST", "/v1/user/devices", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a GET /user/devices", func() {
			req, _ := http.NewRequest("GET", "/v1/user/devices", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a GET /user/messages/history", func() {
			req, _ := http.NewRequest("GET", "/v1/user/messages/history", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a GET /user/messages/sync/:mid", func() {
			req, _ := http.NewRequest("GET", "/v1/user/messages/sync/5", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a DELETE /user/messages/:mid", func() {
			req, _ := http.NewRequest("DELETE", "/v1/user/messages/5", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a POST /user/signout", func() {
			req, _ := http.NewRequest("POST", "/v1/user/signout", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a POST /user/contacts", func() {
			req, _ := http.NewRequest("POST", "/v1/user/contacts", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		g.It("Should allow a GET /user/contacts", func() {
			req, _ := http.NewRequest("GET", "/v1/user/contacts", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	})
}
