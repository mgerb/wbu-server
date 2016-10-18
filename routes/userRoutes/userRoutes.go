package userRoutes

import (
	"log"
	"time"

	"../../config"
	"../../operations/userOperations"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

//HandleTest - test function for random things
func HandleTest(ctx *iris.Context) {
	log.Println(userOperations.CreateUser("username", "password"))
	response := "test"
	ctx.Write(response)
}

//CreateUser - create user account - currently takes in username and password
func CreateUser(ctx *iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")

	err := userOperations.CreateUser(username, password)

	if err == nil {
		ctx.JSON(200, `{"message": "Account Created"}`)
	} else {
		ctx.JSON(500, `{"message": "`+err.Error()+`"}`)
	}
}

//Login - log the user in - on success send jwt
func Login(ctx *iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")

	if userOperations.ValidLogin(username, password) == true {
		id := userOperations.GetUserID(username)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"id":       id,
			"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(config.Config.TokenSecret))
		ctx.JSON(200, `{"jwt": `+tokenString+`}`)

	} else {
		ctx.JSON(500, `{"message": "Invalid login credentials"}`)
	}
}
