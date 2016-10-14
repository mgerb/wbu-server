package routes

import (
	"./middleware"
	"./userRoutes"
	"github.com/kataras/iris"
)

func Routes() *iris.Framework {
	app := iris.New()

	middleware.ApplyMiddleware(app)

	app.Get("/test", userRoutes.HandleTest)
	app.Post("/login", userRoutes.Login)

	return app
}
