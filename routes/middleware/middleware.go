package middleware

import (
	"log"

	"../../config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
)

//ApplyMiddleware - applies middleware to iris framework
func ApplyMiddleware(ctx *iris.Framework) {
	ctx.Use(logger.New())
	ctx.UseFunc(checkJWT)
}

func checkJWT(ctx *iris.Context) {
	path := ctx.PathString()

	if path == "/user/login" || path == "/user/createUser" {
		ctx.Next()
	} else {
		//get Authorization header - jwt token
		authToken := ctx.RequestHeader("Authorization")

		//parse the token
		//TODO - FIX ERROR HANDLING HERE
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.TokenSecret), nil
		})

		//check if actual token
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(500, `{"message": "Invalid Token"}`)
		} else {
			//get the claims from token - userName and id
			if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
				ctx.Set("userName", claims["userName"].(string))
				ctx.Set("userID", claims["userID"].(string))
				ctx.Next()
			} else {
				log.Println(err.Error())
				ctx.JSON(500, `{"message": "Invalid Authentication"}`)
			}
		}
	}
}
