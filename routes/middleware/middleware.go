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
	app.Use(checkJWT)

	if !config.Flags.Production {
		app.Use(middleware.Logger())
	}
}

//define custom JWT middleware
func checkJWT(next echo.HandlerFunc) echo.HandlerFunc {

	//return handler function
	return func(ctx echo.Context) error {
		path := ctx.Request().URL().Path()

		if bypassRoutes(path) {
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
				firstName, ok_firstName := claims["firstName"]
				lastName, ok_lastName := claims["lastName"]

				if ok_email && ok_userID && ok_firstName && ok_lastName {
					ctx.Set("email", email.(string))
					ctx.Set("userID", userID.(string))
					ctx.Set("firstName", firstName.(string))
					ctx.Set("lastName", lastName.(string))
				} else {
					return ctx.JSON(500, response.Json("Token claims error.", response.INVALID_AUTHENTICATION))
				}

				return next(ctx)
			}
			return ctx.JSON(500, response.Json("Invalid authentication.", response.INVALID_AUTHENTICATION))

		default:
			log.Println(err.Error())
			return ctx.JSON(500, response.Json("Invalid token.", response.INVALID_AUTHENTICATION))

		}
	}

}

var prodRoutes = map[string]bool{
	"/user/loginFacebook": true,
}

var devRoutes = map[string]bool{
	"/user/loginFacebook": true,
	"/user/login":         true,
	"/user/createUser":    true,
	"/test":               true,
}

// configure routes to bypass the authentication middleware
func bypassRoutes(path string) bool {

	switch config.Flags.Production {
	case true:
		return prodRoutes[path]
	}

	return devRoutes[path]
}
