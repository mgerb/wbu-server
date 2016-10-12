package operations

import (
	"errors"
	"strconv"

	"../db"
	"../model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Groups []int `json:"groups"` //id of each groups user is in
}

//store username/password in hash
func Create(username string, password string) error {
	if Exists(username) != true {
		new_id, _ := db.Client.Incr(model.USER_KEY_STORE()).Result()

		db.Client.Set(model.USER_NAME(username), new_id, 0)
		db.Client.HMSet(model.USER_HASH(int(new_id)), map[string]string{
			"username": username,
			"password": GenerateHash(password),
		})
		return nil
	}
	return errors.New("username already exists")
}

//check if password and username are correct
func ValidLogin(username string, password string) bool {
	if Exists(username) == true {
		id := GetUserID(username)
		result, _ := db.Client.HGet(model.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(password)) == nil {
			return true
		}
	}
	return false
}

func Exists(username string) bool {
	return db.Client.Get(model.USER_NAME(username)).Err() == nil
}

func GetUserID(username string) int {
	result, _ := db.Client.Get(model.USER_NAME(username)).Result()
	i, _ := strconv.Atoi(result)
	return i
}

func GenerateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
