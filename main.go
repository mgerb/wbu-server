package main

import (
	"./config"
	"./db"
	"./routes"
)

func main() {
	config := config.ReadConfig()
	db.Configure(config)

	routes.Routes().Listen(":" + config.ServerPort)
}
