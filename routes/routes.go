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

	//GET---------------------------------------------------

	//user
	app.GET("/test", userRoutes.HandleTest)
	app.GET("/user/userGroups", userRoutes.GetGroups)
	app.GET("/user/groupInvites", userRoutes.GetInvites)

	//groups
	app.GET("/group/members/:groupID", groupRoutes.GetMembers)
	app.GET("/group/messages/:groupID", groupRoutes.GetMessages)

	//POST---------------------------------------------------
	//user
	app.POST("/user/createUser", userRoutes.CreateUser)
	app.POST("/user/login", userRoutes.Login)

	//groups
	app.POST("/group/createGroup", groupRoutes.CreateGroup)
	app.POST("/group/inviteUser", groupRoutes.InviteToGroup)
	app.POST("/group/storeMessage/:groupID", groupRoutes.StoreMessage)
}
