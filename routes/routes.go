package routes

import "github.com/kataras/iris"

func Routes() *iris.Framework {
	app := iris.New()

	app.Get("/test", handler)
	return app
}

func handler(c *iris.Context) {
	c.Write("Hello test")
}
