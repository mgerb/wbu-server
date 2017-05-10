package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mgerb/wbu-server/config"
)

const (
	directory = "./changescripts"
)

func main() {

	config.ReadConfigFile()

	// open database connection
	sql, err := sql.Open("sqlite3", config.Config.DatabaseName+"?cache=shared&mode=rwc")

	if err != nil {
		log.Fatal(err)
	}

	defer sql.Close()

	sql.SetMaxOpenConns(1)

	// load scripts into memory
	changeScripts := loadScripts()

	// loop through scripts
	for _, script := range changeScripts {
		fmt.Println(script)

		// execute script
		_, err := sql.Exec(script)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func loadScripts() []string {

	files, _ := ioutil.ReadDir(directory)

	var err error
	var fileContents []byte

	scripts := []string{}

	for _, f := range files {

		// skip if not .sql file
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}

		fileContents, err = ioutil.ReadFile(directory + f.Name())

		if err != nil {
			log.Fatal(err)
		}

		scripts = append(scripts, string(fileContents))
	}

	return scripts
}
