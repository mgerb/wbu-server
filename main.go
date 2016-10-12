package main

import (
	"fmt"

	"./config"
	"./db"
	"./operations"
)

func main() {
	config := config.ReadConfig()
	db.Configure(config)

	user := operations.User{
		"Mitchell",
		"password",
	}

	err := user.Create()
	fmt.Println(err)
	fmt.Println(user.ValidLogin())
}
