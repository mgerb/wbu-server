package routes

import (
    "fmt"
    
    "github.com/kataras/iris"
    "../operations"
)

func HandleUser(ctx *iris.Context){
    fmt.Println(operations.Create("username", "password"))
    //response := ctx.RequestHeader("User-Agent") + "\n" + string(ctx.RequestURI())
    response := "test"
    ctx.Write(response)
}