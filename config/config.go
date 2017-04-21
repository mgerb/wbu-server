package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	Config configStruct
	Flags  configFlags
)

type configFlags struct {
	Production bool
}

type configStruct struct {
	Address      string `json:"Address"`
	TokenSecret  string `json:"TokenSecret"`
	DatabaseName string `json:"DatabaseName"`
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
	} else {
		fmt.Println("Starting development mode...")
	}
}
