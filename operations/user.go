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
func (u *User) Create() error {
	if u.Exists() != true {
		new_id, _ := db.Client.Incr(model.USER_KEY_STORE()).Result()

		db.Client.Set(model.USER_NAME(u.Username), new_id, 0)
		db.Client.HMSet(model.USER_HASH(int(new_id)), map[string]string{
			"username": u.Username,
			"password": GenerateHash(u.Password),
		})
		return nil
	}
	return errors.New("username already exists")
}

//check if password and username are correct
func (u *User) ValidLogin() bool {
	if u.Exists() == true {
		id := u.GetUserID()
		result, _ := db.Client.HGet(model.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(u.Password)) == nil {
			return true
		}
	}
	return false
}

func (u *User) Exists() bool {
	return db.Client.Get(model.USER_NAME(u.Username)).Err() == nil
}

func (u *User) GetUserID() int {
	result, _ := db.Client.Get(model.USER_NAME(u.Username)).Result()
	i, _ := strconv.Atoi(result)
	return i
}

func GenerateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
