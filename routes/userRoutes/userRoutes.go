package userRoutes

import (
	"../../operations/groupOperations"
	"../../operations/userOperations"

	"../../utils/response"
	"github.com/labstack/echo"
)

//HandleTest - test function for random things
func HandleTest(ctx echo.Context) error {
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

	return ctx.JSON(200, map[string]string{"message": message})
}

//CreateUser - create user account - currently takes in userName and password
func CreateUser(ctx echo.Context) error {
	userName := ctx.FormValue("userName")
	password := ctx.FormValue("password")

	err := userOperations.CreateUser(userName, password)

	if err == nil {
		return ctx.JSON(200, map[string]string{"message": "Account Created"})
	} else {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}
}

//Login - log the user in - on success send jwt
func Login(ctx echo.Context) error {
	userName := ctx.FormValue("userName")
	password := ctx.FormValue("password")

	jwt, err := userOperations.Login(userName, password)

	if err == nil {
		return ctx.JSON(200, map[string]string{"token": jwt})
	} else {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}
}

func GetGroups(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	groups, err := userOperations.GetGroups(userID)

	if err == nil {
		return ctx.JSON(200, map[string]interface{}{"groups": groups})
	} else {
		return ctx.JSON(500, response.Json("Unable to get groups.", response.INTERNAL_ERROR))
	}
}
