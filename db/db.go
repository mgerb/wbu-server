package db

import (
	"fmt"

	redis "gopkg.in/redis.v4"
)

var Client *redis.Client

func Configure(address string, password string) {
	options := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}

	Client = redis.NewClient(options)

	fmt.Println("Database configured")
}
