package main

import (
	"fmt"

	"./config"
	"./db"
	"./operations"
	"./routes"
)

func main() {
	config := config.ReadConfig()
	db.Configure(config)

	err := operations.Create("mitchell", "gerber")
	fmt.Println(err)
	fmt.Println(operations.ValidLogin("mitchell", "gerber"))

	routes.Routes().Listen(":8080")
}
