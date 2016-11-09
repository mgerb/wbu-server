package db

import (
	"fmt"

	redis "gopkg.in/redis.v5"
)

var Client *redis.Client

func Connect(address string, password string) {

	options := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}

	Client = redis.NewClient(options)

	fmt.Println("Database configured")
}
