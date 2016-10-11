package db

import (
	"fmt"

	"../config"
	"gopkg.in/redis.v4"
)

var Client *redis.Client

func Configure(c config.Config) {
	options := &redis.Options{
		Addr:     c.DatabaseAddress,
		Password: c.DatabasePassword,
		DB:       0,
	}

	Client = redis.NewClient(options)

	fmt.Println("Database configured")
}
