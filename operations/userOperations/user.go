package userOperations

import (
	"errors"
	"strconv"
	"time"
	"regexp"
	"../../config"
	"../../db"
	"../../model/userModel"
	"../../utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

//CreateUser - store username/password in hash
func CreateUser(username string, password string) error {
	
	//DO VALIDATION
	if !regexp.MustCompile(utils.UsernameRegex).MatchString(username){
		return errors.New("Invalid username.")
	}	
	
	//check if the username already exists in redis
	_, err := GetUserID(username)
	if err == nil {
		return errors.New("Username already exists.")
	}

	temp, _ := db.Client.Incr(userModel.USER_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.Set(userModel.USER_ID(username), newID, 0)

	//set user object in redis
	pipe.HMSet(userModel.USER_HASH(newID), map[string]string{
		"username": username,
		"password": generateHash(password),
	})

	_, err = pipe.Exec()

	return err
}

//seconds in 30 days
var expirationTime int64 = 30 * 24 * 60 * 60

//ValidLogin - check if password and username are correct
func Login(username string, password string) (string, error) {
	id, err := GetUserID(username)

	if err == nil {
		result, _ := db.Client.HGet(userModel.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(password)) != nil {
			return "", errors.New("Invalid password.")
		}
	} else {
		return "", errors.New("User does not exist.")
	}

	//if user has valid login - generate jwt	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))
	
	if tokenError != nil {
		return "", errors.New("Token error.")
	}
	
	return tokenString, nil
}

//GetUserID = return user id as string
func GetUserID(username string) (string, error) {
	return db.Client.Get(userModel.USER_ID(username)).Result()
}

//GetUserGroups - get all the groups the user exists in
func GetUserGroups(userID string) ([]string, error) {
	return db.Client.SMembers(userModel.USER_GROUPS(userID)).Result()
}

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}

//
//
//
//TODO
func JoinGroup(userID string) error {
	return errors.New("TODO")
}

func LeaveGroup(userID string, groupid string) error {
	return errors.New("TODO")
}
