package userRoutes

import (
	"github.com/labstack/echo"
	"github.com/mgerb/wbu-server/operations/userOperations"
	"github.com/mgerb/wbu-server/utils/response"
	"github.com/mgerb/wbu-server/utils/tokens"
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

		messages, err := groupOperations.GetUserGroupMessages("1", "fbID:10207835974837361")

		switch err {
		case nil:
			return ctx.JSON(200, messages)
		default:
			return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
		}
	*/
	return nil
}

//CreateUser - create user account - currently takes in userName and password
func CreateUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	firstName := ctx.FormValue("firstName")
	lastName := ctx.FormValue("lastName")

	err := userOperations.CreateUser(email, password, firstName, lastName)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Account created.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// TODO
// DeleteUser - deletes all user information based on their userID
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

// Login - log the user in - on success send jwt
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

// LoginFacebook -
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

// SearchUserByName -
func SearchUserByName(ctx echo.Context) error {
	name := ctx.Param("name")
	userID := ctx.Get("userID").(string)

	userList, err := userOperations.SearchUserByName(name, userID)

	switch err {
	case nil:
		return ctx.JSON(200, userList)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// RefreshJWT -
func RefreshJWT(ctx echo.Context) error {
	email := ctx.Get("email").(string)
	userID := ctx.Get("userID").(string)
	firstName := ctx.Get("firstName").(string)
	lastName := ctx.Get("lastName").(string)

	token, lastRefreshTime, err := tokens.GetJWT(email, userID, firstName, lastName)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"jwt": token, "lastRefreshTime": lastRefreshTime})

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// UpdateFCMToken -
func UpdateFCMToken(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	token := ctx.FormValue("token")

	err := userOperations.UpdateFCMToken(userID, token)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("token updated", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// ToggleNotifications -
func ToggleNotifications(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	toggle := ctx.FormValue("toggle")

	err := userOperations.ToggleNotifications(userID, toggle)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("notifications updated", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// GetUserSettings -
func GetUserSettings(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	settings, err := userOperations.GetUserSettings(userID)

	switch err {
	case nil:
		return ctx.JSON(200, settings)

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// RemoveFCMToken-
func RemoveFCMToken(ctx echo.Context) error {
	token := ctx.FormValue("token")

	err := userOperations.RemoveFCMToken(token)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Token removed.", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// StoreUserFeedback -
func StoreUserFeedback(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	feedback := ctx.FormValue("feedback")

	err := userOperations.StoreUserFeedback(userID, feedback)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Feedback stored.", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
