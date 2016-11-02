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

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Account Created.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//Login - log the user in - on success send jwt
func Login(ctx echo.Context) error {
	userName := ctx.FormValue("userName")
	password := ctx.FormValue("password")

	jwt, err := userOperations.Login(userName, password)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]string{"token": jwt})
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetGroups(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	groups, err := userOperations.GetGroups(userID)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"groups": groups})
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetInvites(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	invites, err := userOperations.GetInvites(userID)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"invites": invites})

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func JoinGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)

	groupID := ctx.FormValue("groupID")
	groupName := ctx.FormValue("groupName")

	err := userOperations.JoinGroup(userID, userName, groupID, groupName)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Joined group.", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
