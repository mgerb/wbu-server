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
	app.GET("/user/refreshJWT", userRoutes.RefreshJWT)

	//groups
	app.GET("/group/members/:groupID", groupRoutes.GetGroupMembers)
	app.GET("/group/messages/:groupID", groupRoutes.GetMessages)
	app.GET("/group/getGeoLocations/:groupID", groupRoutes.GetGeoLocations)

	//user
	app.POST("/user/createUser", userRoutes.CreateUser)
	app.POST("/user/deleteUser", userRoutes.DeleteUser)
	app.POST("/user/login", userRoutes.Login)
	app.POST("/user/loginFacebook", userRoutes.LoginFacebook)
	app.POST("/user/searchByName", userRoutes.SearchUserByName)

	//groups
	app.POST("/group/createGroup", groupRoutes.CreateGroup)
	app.POST("/group/inviteUser", groupRoutes.InviteUser)
	app.POST("/group/joinGroup", groupRoutes.JoinGroup)
	app.POST("/group/leaveGroup", groupRoutes.LeaveGroup)
	app.POST("/group/deleteGroup", groupRoutes.DeleteGroup)
	app.POST("/group/storeMessage", groupRoutes.StoreMessage)
	app.POST("/group/searchPublicGroups", groupRoutes.SearchPublicGroups)

	//geo
	app.POST("/group/storeGeoLocation", groupRoutes.StoreGeoLocation)
}
