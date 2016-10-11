package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ServerPort       string `json:"ServerPort"`
	DatabaseAddress  string `json:"DatabaseAddress"`
	DatabasePassword string `json:"DatabasePassword"`
}

func ReadConfig() Config {

	log.Println("Reading config file...")

	file, e := ioutil.ReadFile("./config.json")

	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	log.Printf("%s\n", string(file))

	var result Config

	err := json.Unmarshal(file, &result)

	if err != nil {
		log.Println(err)
	}

	return result
}
