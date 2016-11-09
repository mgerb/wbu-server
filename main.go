package main

import (
	"./config"
	"./db"
	"./routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
)

func main() {

	//read config files/flags
	config.Init()

	//connect to database
	db.Connect(config.Config.DatabaseAddress, config.Config.DatabasePassword)

	app := echo.New()
	routes.RegisterRoutes(app)

	app.Run(fasthttp.WithConfig(config.Config.ServerConfig))
}
