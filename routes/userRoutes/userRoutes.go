package userRoutes

import (
	"../../operations/groupOperations"
	"../../operations/userOperations"

	"../../utils/response"
	"github.com/kataras/iris"
)

//HandleTest - test function for random things
func HandleTest(ctx *iris.Context) {
	//response := groupOperations.GetGroupMembers("1")
	//res, _ := json.Marshal(response)
	/*
		pipe := db.Client.Pipeline()
		defer pipe.Close()

		for i := 0; i < 1000000; i++ {
			s := strconv.Itoa(i)
			pipe.HMSet("test"+s, map[string]string{"test1": s, "test2": s})
			log.Println("group created " + s)
		}
		pipe.Exec()
	*/
	/*
		err := groupOperations.StoreMessage("groupid", "userid", "userName", "message")
		message := "success"
		if err != nil {
			message = "error"
		}
	*/

	err := groupOperations.StoreGeoLocation("groupID", "test", "13.4", "userID", "userName")

	message := "success"
	if err != nil {
		message = err.Error()
	}

	ctx.JSON(200, `{"message":"`+message+`"}`)
}

//CreateUser - create user account - currently takes in userName and password
func CreateUser(ctx *iris.Context) {
	userName := ctx.PostValue("userName")
	password := ctx.PostValue("password")

	err := userOperations.CreateUser(userName, password)

	if err == nil {
		ctx.JSON(200, `{"message": "Account Created"}`)
	} else {
		ctx.JSON(500, `{"message": "`+err.Error()+`"}`)
	}
}

//Login - log the user in - on success send jwt
func Login(ctx *iris.Context) {
	userName := ctx.PostValue("userName")
	password := ctx.PostValue("password")

	jwt, err := userOperations.Login(userName, password)

	if err == nil {
		ctx.JSON(200, `{"token": `+jwt+`}`)
	} else {
		ctx.JSON(500, `{"error": `+err.Error()+`}`)
	}
}

func GetUserGroups(ctx *iris.Context) {
	userID := ctx.Get("userID").(string)

	groups, err := userOperations.GetUserGroups(userID)

	if err == nil {
		ctx.JSON(200, groups)
	} else {
		ctx.JSON(500, response.Json("Unable to get groups.", response.INTERNAL_ERROR))
	}
}
