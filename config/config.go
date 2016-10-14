package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Config configStruct

type configStruct struct {
	ServerPort       string `json:"ServerPort"`
	DatabaseAddress  string `json:"DatabaseAddress"`
	DatabasePassword string `json:"DatabasePassword"`
	TokenSecret      string `json:"TokenSecret"`
}

func ReadConfig() {

	log.Println("Reading config file...")

	file, e := ioutil.ReadFile("./config.json")

	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	log.Printf("%s\n", string(file))

	err := json.Unmarshal(file, &Config)

	if err != nil {
		log.Println(err)
	}

}
