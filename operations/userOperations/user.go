package userOperations

import (
	"errors"
	"strconv"

	"../../db"
	"../../model/userModel"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser - store username/password in hash
func CreateUser(username string, password string) error {
	if Exists(username) != true {
		temp, _ := db.Client.Incr(userModel.USER_KEY_STORE()).Result()
		newID := strconv.FormatInt(temp, 10)
		db.Client.Set(userModel.USER_ID(username), newID, 0)
		db.Client.HMSet(userModel.USER_HASH(newID), map[string]string{
			"username": username,
			"password": generateHash(password),
		})
		return nil
	}
	return errors.New("User already exists.")
}

//ValidLogin - check if password and username are correct
func ValidLogin(username string, password string) bool {
	if Exists(username) == true {
		id := GetUserID(username)
		result, _ := db.Client.HGet(userModel.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(password)) == nil {
			return true
		}
	}
	return false
}

//Exists - return if user exists in redis
func Exists(username string) bool {
	return db.Client.Get(userModel.USER_ID(username)).Err() == nil
}

//GetUserID = return user id as string
func GetUserID(username string) string {
	result, _ := db.Client.Get(userModel.USER_ID(username)).Result()
	return result
}

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
