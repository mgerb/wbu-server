package model

import (
	"errors"
	"fmt"
	"strconv"

	"../db"
	"../utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Create() error {
	fmt.Println(u.GetUserID())
	if u.Exists() == true {
		u.hashPassword()
		new_id := db.Client.Incr(USER_KEY_STORE)
		_ = db.Client.Set(USER_NAME+u.Username, new_id.Val(), 0)
		_ = db.Client.HMSet(USER+strconv.FormatInt(new_id.Val(), 10), utils.StructToMap(u))
		return nil
	}

	return errors.New("username already exists")
}

/*
func (u *User) CheckLogin() error {
	if u.Exists() == true {

	}
	return errors.New("Username does not exists")
}
*/

func (u *User) Exists() bool {
	return db.Client.Get(USER_NAME+u.Username).Err() != nil
}

func (u *User) GetUserID() string {
	return db.Client.Get(USER_NAME + u.Username).Val()
}

/*
func (u *User) GetUserID() *redis.StringCmd {
	return db.Client.Get(USER_NAME + u.Username)
}
*/

func (u *User) hashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	u.Password = string(hash)
}
