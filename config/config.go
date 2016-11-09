package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/labstack/echo/engine"
)

var Config configStruct
var Flags configFlags

type configFlags struct {
	Production bool
}

type configStruct struct {
	ServerConfig     engine.Config `json:"ServerConfig"`
	DatabaseAddress  string        `json:"DatabaseAddress"`
	DatabasePassword string        `json:"DatabasePassword"`
	TokenSecret      string        `json:"TokenSecret"`
}

func Init() {
	readConfigFile()
	parseFlags()
}

func readConfigFile() {

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

func parseFlags() {
	Flags.Production = false

	boolPtr := flag.Bool("p", false, "Production mode")

	flag.Parse()

	if *boolPtr {
		fmt.Println("Starting production mode...")
		Flags.Production = true
	}
}
