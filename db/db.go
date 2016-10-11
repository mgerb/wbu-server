package db

import (
    "fmt"
    "gopkg.in/redis.v4"
)

var Client *redis.Client


func Configure(){
    options := &redis.Options{
        Addr: "45.55.230.86:6279",
        Password: "thereisnospoon",
        DB: 0,
    }
    
    Client = redis.NewClient(options)
    
    fmt.Println("Database configured")
}