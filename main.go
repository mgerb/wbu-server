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

	//connect to database and set up client
	db.Connect(config.Config.DatabaseAddress, config.Config.DatabasePassword)
	defer db.SQL.Close()

	app := echo.New()
	routes.RegisterRoutes(app)

	app.Run(fasthttp.WithConfig(config.Config.ServerConfig))
}
