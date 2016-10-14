package userRoutes

import (
	"fmt"
	"time"

	"../../config"
	"../../operations"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

func HandleTest(ctx *iris.Context) {
	fmt.Println(operations.Create("username", "password"))
	//response := ctx.RequestHeader("User-Agent") + "\n" + string(ctx.RequestURI())
	response := "test"
	ctx.Write(response)
}

func Login(ctx *iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")
	if operations.ValidLogin(username, password) == true {
		id := operations.GetUserID(username)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"id":       id,
			"foo":      "bar",
			"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(config.Config.TokenSecret))
		ctx.JSON(200, `{"jwt": `+tokenString+`}`)
	} else {
		ctx.JSON(500, `{"message": "Invalid login credentials"}`)
	}
}

func CreateUser(ctx *iris.Context) {

}
