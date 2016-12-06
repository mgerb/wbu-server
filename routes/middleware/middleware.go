package middleware

import (
	"log"

	"../../config"
	"../../utils/response"
	"github.com/dgrijalva/jwt-go"
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

//TODO - jwt refreshing
//define custom JWT middleware
func checkJWT(next echo.HandlerFunc) echo.HandlerFunc {

	//return handler function
	return func(ctx echo.Context) error {
		path := ctx.Request().URL().Path()

		//routes to skip authentication
		switch path {
		case "/user/loginFacebook",
			"/user/login",
			"/user/createUser",
			"/test":
			return next(ctx)
		}

		//get Authorization header - jwt token
		authToken := ctx.Request().Header().Get("Authorization")

		//parse the token
		//TODO - FIX ERROR HANDLING HERE
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.TokenSecret), nil
		})

		switch err {
		case nil:
			if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {

				email, ok_email := claims["email"]
				userID, ok_userID := claims["userID"]
				fullName, ok_fullName := claims["fullName"]

				if ok_email && ok_userID && ok_fullName {
					ctx.Set("email", email.(string))
					ctx.Set("userID", userID.(string))
					ctx.Set("fullName", fullName.(string))
				} else {
					return ctx.JSON(500, response.Json("Token claims error.", response.INTERNAL_ERROR))
				}

				return next(ctx)
			}
			return ctx.JSON(500, response.Json("Invalid authentication.", response.INTERNAL_ERROR))

		default:
			log.Println(err.Error())
			return ctx.JSON(500, response.Json("Invalid token.", response.INTERNAL_ERROR))

		}
	}

}
