package main

import (
	"log"
	"net/http"
	"os"
	"portal-server/api/controller/access"
	"portal-server/api/controller/user"
	"portal-server/api/middleware"
	"portal-server/api/util"
	"portal-server/store"

	"github.com/gin-gonic/gin"
)

var (
	dbName     = os.Getenv("DB_NAME")
	dbUser     = os.Getenv("DB_API_USER")
	dbPassword = os.Getenv("DB_API_PASSWORD")
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
			secure.GET("/messages/sync/:mid", user.SyncMessagesEndpoint)
			secure.DELETE("/messages/:mid", user.DeleteMessageEndpoint)
			secure.POST("/contacts")
			secure.GET("/contacts")
			secure.DELETE("/contacts/:cid")
			secure.POST("/signout", user.SignoutEndpoint)
		}
	}
	return r
}

func main() {
	if util.GcmApiKey == "" || util.GcmSenderID == "" {
		log.Fatalln("Missing GCM_SENDER_ID or GCM_API_KEY environment variables")
	}

	if dbName == "" || dbUser == "" || dbPassword == "" {
		log.Fatalln("Missing DB_NAME, DB_API_USER, or DB_API_PASSWORD environment variables")
	}

	store := store.GetStore(dbName, dbUser, dbPassword)
	httpClient := http.DefaultClient
	API(store, httpClient).Run(":8080")
}
