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

	//get requests

	//user
	app.Get("/test", userRoutes.HandleTest)
	app.Get("/userGroups", userRoutes.GetUserGroups)

	//groups

	//post requests
	//user
	app.Post("/createUser", userRoutes.CreateUser)
	app.Post("/login", userRoutes.Login)

	//groups
	app.Post("/createGroup", groupRoutes.CreateGroup)
	app.Post("/group/storeMessage", groupRoutes.StoreMessage)

	return app
}
