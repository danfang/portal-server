package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
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

		g.It("Should allow a GET /user/message/history", func() {
			req, _ := http.NewRequest("GET", "/v1/user/messages/history", nil)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	})
}
