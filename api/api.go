// The Portal Messaging API.
//
//     Schemes: http
//     Host: 52.89.157.164:8080
//     BasePath: /v1
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"net/http"
	"portal-server/api/auth"
	"portal-server/api/routing/access"
	"portal-server/api/routing/user"
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
	r.Use(auth.CORSMiddleware())

	// Add swagger.json file
	r.StaticFile("/swagger.json", "./api/swagger.json")

	v1 := r.Group("/v1")
	{
		accessGroup := v1.Group("/")
		{
			accessRouter := access.Router{db, httpClient}

			// swagger:route POST /register register
			//
			// Register a user via email and password.
			//
			//     Consumes:
			//     - application/json
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			// 		 default: detailError
			//       200: success
			// 		 400: error
			//		 500: error
			accessGroup.POST("/register", accessRouter.RegisterEndpoint)

			// swagger:route POST /login login
			//
			// User login via email and password.
			//
			//     Consumes:
			//     - application/json
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			// 		 default: detailError
			//       200: loginResponse
			// 		 400: error
			//		 500: error
			accessGroup.POST("/login", accessRouter.LoginEndpoint)

			// swagger:route POST /login/google googleLogin
			//
			// Login or register via a Google account.
			//
			//     Consumes:
			//     - application/json
			//
			//     Produces:
			//     - application/json
			//
			//     Responses:
			// 	     default: detailError
			//       200: loginResponse
			// 	     400: error
			//	     500: error
			accessGroup.POST("/login/google", accessRouter.GoogleLoginEndpoint)

			// swagger:route GET /verify/{token} verifyToken
			//
			// Consume a user email verification token.
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			//       200: success
			// 	     400: error
			//		 500: error
			accessGroup.GET("/verify/:token", accessRouter.VerifyUserEndpoint)
		}

		userGroup := v1.Group("/user")
		userGroup.Use(auth.AuthenticationMiddleware(db))
		{
			userRouter := user.Router{db, httpClient}

			// swagger:route POST /user/devices user addDevice
			//
			// Register a new Google Cloud Messaging device.
			//
			//     Consumes:
			//     - application/json
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			//       200: addDevice
			//		 401: error
			//		 500: error
			//       default: detailError
			userGroup.POST("/devices", userRouter.AddDeviceEndpoint)

			// swagger:route GET /user/devices user getDevices
			//
			// Retrieve a user's existing connected devices.
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			//       200: deviceList
			//       401: error
			//       500: error
			//       default: detailError
			userGroup.GET("/devices", userRouter.GetDevicesEndpoint)

			// swagger:route GET /user/messages/history user messageHistory
			//
			// Retrieve a user's existing connected devices.
			//
			//     Produces:
			//     - application/json
			//
			//     Schemes: http
			//
			//     Responses:
			//       200: messageHistory
			//       401: error
			//       500: error
			//       default: detailError
			userGroup.GET("/messages/history", userRouter.GetMessageHistoryEndpoint)
		}
	}
	return r
}

func main() {
	db := model.GetDB(dbUser, dbName, dbPassword)
	API(db).Run(":8080")
}
