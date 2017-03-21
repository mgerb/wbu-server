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

	db.Start(config.Config.DatabaseName)
	defer db.SQL.Close()

	app := echo.New()
	routes.RegisterRoutes(app)

	app.Run(fasthttp.WithConfig(config.Config.ServerConfig))
}
