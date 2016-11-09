package tokens

import (
	"../../config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	//expiration time for JWT
	expirationTime int64 = 60 * 24 * 60 * 60 //60 days
)

func GetJWT(email string, userID string, fullName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    email,
		"userID":   userID,
		"fullName": fullName,
		"exp":      time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))

	if tokenError != nil {
		return "", errors.New("Token error.")
	}

	return tokenString, nil
}
