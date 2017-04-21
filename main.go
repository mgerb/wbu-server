package main

import (
	"./config"
	"./db"
	"./routes"

	"golang.org/x/crypto/acme/autocert"
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
		//app.AutoTLSManager.HostPolicy = autocert.HostWhitelist("redis.mitchellgerber.com")
		app.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		app.Logger.Fatal(app.StartAutoTLS("0.0.0.0:443"))
	} else {
		app.Logger.Fatal(app.Start(config.Config.Address))
	}
}
