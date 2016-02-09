package main

import (
	"net/http"
	"portal-server/api/controller/access"
	"portal-server/api/controller/user"
	"portal-server/api/middleware"
	"portal-server/store"

	"github.com/gin-gonic/gin"
)

const (
	dbUser     = "portal_api"
	dbPassword = "password"
)

// API returns a Gin router based on a given database.
func API(store store.Store, httpClient *http.Client) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Set context variables
	r.Use(middleware.SetStore(store))
	r.Use(middleware.SetWebClient(httpClient))

	// Add swagger.json file
	r.StaticFile("/swagger.json", "./api/swagger.json")

	v1 := r.Group("/v1")
	{
		base := v1.Group("/")
		{
			base.POST("/register", access.RegisterEndpoint)
			base.POST("/login", access.LoginEndpoint)
			base.POST("/login/google", access.GoogleLoginEndpoint)
			base.GET("/verify/:token", access.VerifyUserEndpoint)
		}

		secure := v1.Group("/user")
		secure.Use(middleware.AuthenticationMiddleware())
		{
			secure.POST("/devices", user.AddDeviceEndpoint)
			secure.GET("/devices", user.GetDevicesEndpoint)
			secure.GET("/messages/history", user.GetMessageHistoryEndpoint)
			secure.POST("/signout", user.SignoutEndpoint)
		}
	}
	return r
}

func main() {
	store := store.GetStore(dbUser, dbPassword)
	httpClient := http.DefaultClient
	API(store, httpClient).Run(":8080")
}
