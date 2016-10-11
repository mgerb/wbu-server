package main

import(
    "./db"
    "./model"
    "fmt"
)

func main(){
    db.Configure()
    
    user := model.User{
        "Mitchell",
        "Password",
    }
    
    err := user.CreateUser()
    fmt.Println(err)
}
