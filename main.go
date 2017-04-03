package main

import (
	"./config"
	"./db"
	"./routes"

	"github.com/labstack/echo"
)

func main() {

	//read config files/flags
	config.Init()

	db.Start(config.Config.DatabaseName)
	defer db.SQL.Close()

	app := echo.New()
	routes.RegisterRoutes(app)

	if config.Flags.Production {
		app.Logger.Fatal(app.StartTLS(config.Config.Address, config.Config.CertFile, config.Config.KeyFile))
	} else {
		app.Logger.Fatal(app.Start(config.Config.Address))
	}
}
