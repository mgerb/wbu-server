package lua

import (
	"fmt"
	"io/ioutil"
	"log"
)

var file = map[string]string{}

//Init - reads all files into a map of strings
func Init() {
	dir := "./luascripts/"

	files, _ := ioutil.ReadDir(dir)

	var err error
	var fileContents []byte

	for _, f := range files {
		fileContents, err = ioutil.ReadFile(dir + f.Name())

		if err != nil {
			log.Println(err)
		}

		file[f.Name()] = string(fileContents)
	}

	fmt.Println("Lua scripts loaded into memory...")
}

//Use - return lua script to use
func Use(fileName string) string {
	return file[fileName]
}
