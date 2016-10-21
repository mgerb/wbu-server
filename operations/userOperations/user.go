package userOperations

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"../../config"
	"../../db"
	"../../model/userModel"
	"../../utils/regex"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser - store userName/password in hash
func CreateUser(userName string, password string) error {

	//DO VALIDATION
	if !regexp.MustCompile(regex.USERNAME).MatchString(userName) {
		return errors.New("Invalid userName.")
	}

	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("Invalid password.")
	}

	//check if the userName already exists in redis
	_, err := GetUserID(userName)
	if err == nil {
		return errors.New("userName already exists.")
	}

	temp, _ := db.Client.Incr(userModel.USER_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.Set(userModel.USER_ID(userName), newID, 0)

	//set user object in redis
	pipe.HMSet(userModel.USER_HASH(newID), map[string]string{
		"userName": userName,
		"password": generateHash(password),
	})

	_, err = pipe.Exec()

	return err
}

//seconds in 30 days
var expirationTime int64 = 30 * 24 * 60 * 60

//ValidLogin - check if password and userName are correct
func Login(userName string, password string) (string, error) {
	id, err := GetUserID(userName)

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
		"userName": userName,
		"userID":   id,
		"exp":      time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))

	if tokenError != nil {
		return "", errors.New("Token error.")
	}

	return tokenString, nil
}

//GetUserID = return user id as string
func GetUserID(userName string) (string, error) {
	return db.Client.Get(userModel.USER_ID(userName)).Result()
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
