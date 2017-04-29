package main

import (
	"github.com/mgerb/wbu-server/config"
	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/routes"

	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	//read config files/flags
	config.Init()

	db.Start(config.Config.DatabaseName)
	defer db.SQL.Close()

	app := echo.New()
	app.Server.WriteTimeout = time.Second * 10
	app.Server.ReadTimeout = time.Second * 10

	routes.RegisterRoutes(app)

	if config.Flags.Production {
		//app.AutoTLSManager.HostPolicy = autocert.HostWhitelist("redis.mitchellgerber.com")
		app.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		app.Logger.Fatal(app.StartAutoTLS("0.0.0.0:443"))
	} else {
		app.Logger.Fatal(app.Start(config.Config.Address))
	}
}
