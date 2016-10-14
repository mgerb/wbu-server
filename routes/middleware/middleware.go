package middleware

import(
    "github.com/kataras/iris"
    //"github.com/iris-contrib/middleware/logger"
    "fmt"
)


type MyMiddleware struct{}

func ApplyMiddleware(ctx *iris.Framework){
    ctx.UseFunc(test)
    ctx.UseFunc(login)
    //ctx.Use(logger.New())
}

func test(ctx *iris.Context){
    fmt.Println(ctx.PathString())
    ctx.Next()
}

func login(ctx *iris.Context){
    if ctx.PathString() == "/login"{
        ctx.Write("you are not logged in")
    } else {
        ctx.Next()
    }
}