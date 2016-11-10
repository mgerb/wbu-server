package routes

import (
	"./groupRoutes"
	"./middleware"
	"./userRoutes"

	"github.com/labstack/echo"
)

//register routes
func RegisterRoutes(app *echo.Echo) {

	middleware.ApplyMiddleware(app)

	app.GET("/test", userRoutes.HandleTest)

	//user
	app.GET("/user/groups", userRoutes.GetGroups)
	app.GET("/user/invites", userRoutes.GetInvites)

	//groups
	app.GET("/group/members/:groupID", groupRoutes.GetMembers)
	app.GET("/group/messages/:groupID", groupRoutes.GetMessages)

	//user
	app.POST("/user/createUser", userRoutes.CreateUser)
	app.POST("/user/login", userRoutes.Login)
	app.POST("/user/loginFacebook", userRoutes.LoginFacebook)

	//groups
	app.POST("/group/createGroup", groupRoutes.CreateGroup)
	app.POST("/group/inviteUser", groupRoutes.InviteUser)
	app.POST("/group/joinGroup", groupRoutes.JoinGroup)
	app.POST("/group/storeMessage", groupRoutes.StoreMessage)
}
