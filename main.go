package main

import (
	"fmt"

	"./config"
	"./db"
	"./model"
)

func main() {
	config := config.ReadConfig()
	db.Configure(config)

	user := model.User{
		"Mitchell",
		"Password",
	}

	err := user.Create()
	fmt.Println(err)
}
