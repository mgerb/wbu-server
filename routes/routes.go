package routes

import (
	"./groupRoutes"
	"./middleware"
	"./userRoutes"
	"github.com/kataras/iris"
)

//register routes
func Routes() *iris.Framework {
	app := iris.New()

	middleware.ApplyMiddleware(app)

	//GET---------------------------------------------------

	//user
	app.Get("/test", userRoutes.HandleTest)
	app.Get("/user/userGroups", userRoutes.GetGroups)

	//groups
	app.Get("/group/members/:groupID", groupRoutes.GetMembers)
	app.Get("/group/messages/:groupID", groupRoutes.GetMessages)

	//POST---------------------------------------------------
	//user
	app.Post("/user/createUser", userRoutes.CreateUser)
	app.Post("/user/login", userRoutes.Login)

	//groups
	app.Post("/group/createGroup", groupRoutes.CreateGroup)
	app.Post("/group/storeMessage/:groupID", groupRoutes.StoreMessage)

	return app
}
