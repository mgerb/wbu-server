package main

import (
	"./config"
	"./db"
	"./routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
)

func main() {
	config.ParseFlags()
	config.ReadConfig()

	db.Configure(config.Config.DatabaseAddress, config.Config.DatabasePassword)

	app := echo.New()
	routes.RegisterRoutes(app)

	app.Run(fasthttp.WithConfig(config.Config.ServerConfig))
}
