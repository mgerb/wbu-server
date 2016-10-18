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

	//groups

	//post requests
	//user
	app.Post("/createuser", userRoutes.CreateUser)
	app.Post("/login", userRoutes.Login)

	//groups
	app.Post("/creategroup", groupRoutes.CreateGroup)

	return app
}
