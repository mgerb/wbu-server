package lua

import (
	"io/ioutil"
	"log"
)

var file = map[string]string{}

//Init - reads all files in to a map of strings
func Init() {
	dir := "./operations/lua/scripts/"

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

}

//Use - return lua script to use
func Use(fileName string) string {
	return file[fileName]
}
