package middleware

import (
	"log"

	"../../config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//ApplyMiddleware - applies middleware to iris framework
func ApplyMiddleware(app *echo.Echo) {
	//app.Use(logger.New())
	app.Use(checkJWT)

	if !config.Flags.Production {
		app.Use(middleware.Logger())
	}
}

//define custom JWT middleware
func checkJWT(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		path := ctx.Request().URL().Path()

		if path == "/user/login" || path == "/user/createUser" {
			return next(ctx)
		} else {
			//get Authorization header - jwt token
			authToken := ctx.Request().Header().Get("Authorization")

			//parse the token
			//TODO - FIX ERROR HANDLING HERE
			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Config.TokenSecret), nil
			})

			//check if actual token
			if err != nil {
				log.Println(err.Error())
				return ctx.JSON(500, map[string]string{"message": "Invalid Token"})
			} else {
				//get the claims from token - userName and id
				if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
					ctx.Set("userName", claims["userName"].(string))
					ctx.Set("userID", claims["userID"].(string))
					return next(ctx)
				} else {
					log.Println(err.Error())
					return ctx.JSON(500, map[string]string{"message": "Invalid Authentication"})
				}
			}
		}
	}
}
