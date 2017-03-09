package tokens

import (
	"errors"
	"time"

	"../../config"
	"github.com/dgrijalva/jwt-go"
)

const (
	//expiration time for JWT
	expirationTime int64 = 60 * 24 * 60 * 60 //60 days
)

func GetJWT(email string, userID string, firstName string, lastName string) (string, int64, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":           email,
		"userID":          userID,
		"firstName":       firstName,
		"lastName":        lastName,
		"lastRefreshTime": time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))

	lastRefreshTime := time.Now().Unix()

	if tokenError != nil {
		return "", lastRefreshTime, errors.New("Token error.")
	}

	return tokenString, time.Now().Unix(), nil
}
