package userRoutes

import (
	"../../operations/groupOperations"
	"../../operations/userOperations"
	"../../utils/response"
	"../../utils/tokens"
	"github.com/labstack/echo"
)

//HandleTest - test function for random things
func HandleTest(ctx echo.Context) error {
	/*
		err := groupOperations.StoreGeoLocation("groupID", "test", "13.4", "userID", "userName")

		message := "success"
		if err != nil {
			message = err.Error()
		}

		return ctx.JSON(200, map[string]string{"message": message})
	*/

	/*
		err := groupOperations.StoreUserGroupMessages("1", "fbID:10207835974837361", "te;lsakjfpo84owjofijsakjfhdasouhrfouashfst123")

		switch err {
		case nil:
			return ctx.JSON(200, response.Json("Message Stored.", response.SUCCESS))
		default:
			return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
		}

	*/
	messages, err := groupOperations.GetUserGroupMessages("1", "fbID:10207835974837361")

	switch err {
	case nil:
		return ctx.JSON(200, messages)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//CreateUser - create user account - currently takes in userName and password
func CreateUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	fullName := ctx.FormValue("fullName")

	err := userOperations.CreateUser(email, password, fullName)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Account created.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//DeleteUser - deletes all user information based on their userID
func DeleteUser(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	err := userOperations.DeleteUser(userID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Account deleted.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//Login - log the user in - on success send jwt
func Login(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	userInfo, err := userOperations.Login(email, password)

	switch err {
	case nil:
		return ctx.JSON(200, userInfo)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func LoginFacebook(ctx echo.Context) error {
	//facebook access token
	accessToken := ctx.FormValue("accessToken")

	//create new jwt for user authentication to this server
	userInfo, err := userOperations.LoginFacebook(accessToken)

	switch err {
	case nil:
		return ctx.JSON(200, userInfo)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetGroups(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	groups, err := userOperations.GetGroups(userID)

	switch err {
	case nil:
		return ctx.JSON(200, groups)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetInvites(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	invites, err := userOperations.GetInvites(userID)

	switch err {
	case nil:
		return ctx.JSON(200, invites)

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func RefreshJWT(ctx echo.Context) error {
	email := ctx.Get("email").(string)
	userID := ctx.Get("userID").(string)
	fullName := ctx.Get("fullName").(string)

	token, lastRefreshTime, err := tokens.GetJWT(email, userID, fullName)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"jwt": token, "lastRefreshTime": lastRefreshTime})

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
