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
	app.GET("/group/getGroupInvites", groupRoutes.GetGroupInvites)

	//user
	app.POST("/user/createUser", userRoutes.CreateUser)
	app.POST("/user/login", userRoutes.Login)
	app.POST("/user/loginFacebook", userRoutes.LoginFacebook)
	app.POST("/user/searchByName", userRoutes.SearchUserByName)
	app.POST("/user/updateFCMToken", userRoutes.UpdateFCMToken)

	//groups
	app.POST("/group/createGroup", groupRoutes.CreateGroup)
	app.POST("/group/searchPublicGroups", groupRoutes.SearchPublicGroups)
	app.POST("/group/joinPublicGroup", groupRoutes.JoinPublicGroup)
	app.POST("/group/inviteUserToGroup", groupRoutes.InviteUserToGroup)
	app.POST("/group/joinGroupFromInvite", groupRoutes.JoinGroupFromInvite)
	app.POST("/group/leaveGroup", groupRoutes.LeaveGroup)

	// messages
	app.GET("/group/getMessages/:groupID/:timestamp", groupRoutes.GetMessages)
	app.POST("/group/storeMessage", groupRoutes.StoreMessage)
}
