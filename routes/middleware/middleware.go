package middleware

import (
	"fmt"

	"../../config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
)

func ApplyMiddleware(ctx *iris.Framework) {
	//ctx.UseFunc(test)
	//ctx.UseFunc(login)
	ctx.Use(logger.New())
	ctx.UseFunc(checkJWT)
}

func test(ctx *iris.Context) {
	fmt.Println(ctx.PathString())
	ctx.Next()
}

func login(ctx *iris.Context) {
	/*
	   if ctx.PathString() == "/login"{
	       ctx.Write("you are not logged in")
	   } else {
	       ctx.Next()
	   }
	*/
	ctx.Next()
}

func checkJWT(ctx *iris.Context) {
	authToken := ctx.RequestHeader("Authorization")

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenSecret), nil
	})

	if err != nil {
		ctx.Write(err.Error())
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
			fmt.Println(claims["username"], claims["id"])
			ctx.Next()
		} else {
			ctx.Write(err.Error())
		}
	}

}
