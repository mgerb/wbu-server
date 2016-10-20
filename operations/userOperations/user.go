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

//ValidLogin - check if password and username are correct
func ValidLogin(username string, password string) bool {
	id, err := GetUserID(username)

	if err == nil {
		result, _ := db.Client.HGet(userModel.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(password)) == nil {
			return true
		}
	}

	return false
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
