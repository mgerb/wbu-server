package groupRoutes

import (
	"fmt"

	"../../operations"
	"github.com/kataras/iris"
)

func HandleUser(ctx *iris.Context) {
	fmt.Println(operations.Create("username", "password"))
	//response := ctx.RequestHeader("User-Agent") + "\n" + string(ctx.RequestURI())
	response := "test"
	ctx.Write(response)
}
