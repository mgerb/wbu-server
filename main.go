package main

import (
	"./config"
	"./db"
	"./routes"
)

func main() {
	config.ReadConfig()
	db.Configure(config.Config.DatabaseAddress, config.Config.DatabasePassword)

	routes.Routes().Listen(":" + config.Config.ServerPort)
}
