package routes

import (
	"./groupRoutes"
	"./middleware"
	"./userRoutes"

	"github.com/labstack/echo"
)

//RegisterRoutes -
func RegisterRoutes(app *echo.Echo) {

	middleware.ApplyMiddleware(app)

	app.GET("/test", userRoutes.HandleTest)

	//user
	app.GET("/user/refreshJWT", userRoutes.RefreshJWT)

	//groups
	app.GET("/group/getUserGroups", groupRoutes.GetUserGroups)
	app.GET("/group/getGroupUsers/:groupID", groupRoutes.GetGroupUsers)

	//user
	app.POST("/user/createUser", userRoutes.CreateUser)
	app.POST("/user/login", userRoutes.Login)
	app.POST("/user/loginFacebook", userRoutes.LoginFacebook)
	app.POST("/user/searchByName", userRoutes.SearchUserByName)

	//groups
	app.POST("/group/createGroup", groupRoutes.CreateGroup)
	app.POST("/group/searchPublicGroups", groupRoutes.SearchPublicGroups)
	app.POST("/group/joinPublicGroup", groupRoutes.JoinPublicGroup)

	// messages
	app.GET("/group/getMessages/:groupID/:timestamp", groupRoutes.GetMessages)
	app.POST("/group/storeMessage", groupRoutes.StoreMessage)
}
