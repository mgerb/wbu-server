package middleware

import (
	"log"
	"net"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mgerb/wbu-server/config"
	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/db/lua"
	"github.com/mgerb/wbu-server/model"
	"github.com/mgerb/wbu-server/utils/response"
)

//ApplyMiddleware -
func ApplyMiddleware(app *echo.Echo) {
	app.Use(rateLimit)
	app.Use(checkJWT)

	if !config.Flags.Production {
		app.Use(middleware.Logger())
	}
}

//define custom JWT middleware
func checkJWT(next echo.HandlerFunc) echo.HandlerFunc {

	//return handler function
	return func(ctx echo.Context) error {
		path := ctx.Request().URL.Path

		if bypassRoutes(path) {
			return next(ctx)
		}

		//get Authorization header - jwt token
		authToken := ctx.Request().Header.Get("Authorization")

		//parse the token
		//TODO - FIX ERROR HANDLING HERE
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.TokenSecret), nil
		})

		switch err {
		case nil:
			if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {

				email, okEmail := claims["email"]
				userID, okUserID := claims["userID"]
				firstName, okFirstName := claims["firstName"]
				lastName, okLastName := claims["lastName"]

				if okEmail && okUserID && okFirstName && okLastName {
					ctx.Set("email", email.(string))
					ctx.Set("userID", userID.(string))
					ctx.Set("firstName", firstName.(string))
					ctx.Set("lastName", lastName.(string))
				} else {
					return ctx.JSON(401, response.Json("Token claims error.", response.INVALID_AUTHENTICATION))
				}

				return next(ctx)
			}
			return ctx.JSON(401, response.Json("Invalid authentication.", response.INVALID_AUTHENTICATION))

		default:
			log.Println(err.Error())
			return ctx.JSON(401, response.Json("Invalid token.", response.INVALID_AUTHENTICATION))

		}
	}

}

var prodRoutes = map[string]bool{
	"/user/loginFacebook": true,
	"/user/login":         true,
	"/user/createUser":    true,
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

// rate limit middleware
func rateLimit(next echo.HandlerFunc) echo.HandlerFunc {

	//return handler function
	return func(ctx echo.Context) error {

		// use lua script
		script := redis.NewScript(lua.Use("RateLimit.lua"))

		// get ip address from ip:port
		ip, _, err := net.SplitHostPort(ctx.Request().RemoteAddr)

		if err != nil {
			return ctx.JSON(500, response.Json("Request error.", response.INTERNAL_ERROR))
		}

		_, err = script.Run(db.RClient, []string{model.RateLimitKey + ip}).Result()

		if err != nil {
			return ctx.JSON(429, response.Json(err.Error(), response.INTERNAL_ERROR))
		}

		return next(ctx)
	}
}
