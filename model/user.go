package model

import (
    //"gopkg.in/redis.v4"
    "../db"
    "golang.org/x/crypto/bcrypt"
    "errors"
    "strconv"
)

type User struct{
    Username string `json:"username"`    
    Password string `json:"password"`
}

func (u *User) CreateUser() error{
    
    if u.Exists() != true{
        id := db.Client.Incr(USER_KEY_STORE)
        
        _ = db.Client.Set(USER_NAME + u.Username, id.Val(), 0)
        
        user := map[string]string{"username": u.Username, "password": string(u.GeneratePasswordHash())}
        
        _ = db.Client.HMSet(USER + strconv.FormatInt(id.Val(), 10), user)
        
        return nil    
    }
    
    return errors.New("username already exists")
}

func (u *User) Exists() bool{
    return db.Client.Get(USER_NAME + u.Username).Err() == nil
}

func (u *User) GeneratePasswordHash() []byte{
    hash,_ := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
    return hash
}