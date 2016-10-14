package routes

import (
	"github.com/kataras/iris"
	"./middleware"
)

func Routes() *iris.Framework {
	app := iris.New()
	
	//app.Use(&middleware.MyMiddleware{})
	//app.Use(logger.New())	
	middleware.ApplyMiddleware(app)
	
	app.Get("/test", HandleUser)
	
	return app
}
