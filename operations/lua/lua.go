package lua

import (
	"fmt"
	"io/ioutil"
)

var File map[string]string

func Init() {

	files, _ := ioutil.ReadDir("./operations/lua/scripts")
	for _, f := range files {
		fmt.Println(f.Name())
	}

}
