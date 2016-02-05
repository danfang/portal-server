package main

import (
	"net/http"
	"portal-server/api/controller/access"
	"portal-server/api/controller/user"
	"portal-server/api/middleware"
	"portal-server/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	dbUser     = "portal_api"
	dbName     = "portal"
	dbPassword = "password"
)

// API returns a Gin router based on a given database.
func API(db *gorm.DB) *gin.Engine {
	httpClient := http.DefaultClient

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Add swagger.json file
	r.StaticFile("/swagger.json", "./api/swagger.json")

	v1 := r.Group("/v1")
	{
		accessGroup := v1.Group("/")
		{
			accessRouter := access.Router{Db: db, HTTPClient: httpClient}
			accessGroup.POST("/register", accessRouter.RegisterEndpoint)
			accessGroup.POST("/login", accessRouter.LoginEndpoint)
			accessGroup.POST("/login/google", accessRouter.GoogleLoginEndpoint)
			accessGroup.GET("/verify/:token", accessRouter.VerifyUserEndpoint)
		}

		userGroup := v1.Group("/user")
		userGroup.Use(middleware.AuthenticationMiddleware(db))
		{
			userRouter := user.Router{Db: db, HTTPClient: httpClient}
			userGroup.POST("/devices", userRouter.AddDeviceEndpoint)
			userGroup.GET("/devices", userRouter.GetDevicesEndpoint)
			userGroup.GET("/messages/history", userRouter.GetMessageHistoryEndpoint)
		}
	}
	return r
}

func main() {
	db := model.GetDB(dbUser, dbName, dbPassword)
	API(db).Run(":8080")
}
