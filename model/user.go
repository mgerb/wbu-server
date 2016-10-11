package model

import (
	"errors"
	"strconv"

	"../db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Create() error {
	if u.Exists() != true {
		new_id, _ := db.Client.Incr(USER_KEY_STORE).Result()

		db.Client.Set(USER_NAME+u.Username, new_id, 0)
		db.Client.HMSet(USER+strconv.FormatInt(new_id, 10), map[string]string{
			"username": u.Username,
			"password": GenerateHash(u.Password),
		})
		return nil
	}
	return errors.New("username already exists")
}

func (u *User) ValidLogin() bool {
	if u.Exists() == true {
		id := u.GetUserID()
		result, _ := db.Client.HGet(USER+id, "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(u.Password)) == nil {
			return true
		}
	}
	return false
}

func (u *User) Exists() bool {
	return db.Client.Get(USER_NAME+u.Username).Err() == nil
}

func (u *User) GetUserID() string {
	id, _ := db.Client.Get(USER_NAME + u.Username).Result()
	return id
}

func GenerateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
