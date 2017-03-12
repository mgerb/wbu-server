package groupRoutes

import (
	"../../operations/groupOperations"
	"../../utils/response"
	"github.com/labstack/echo"
)

//CreateGroup - create a new group with groupName and user id as the owner
func CreateGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupName := ctx.FormValue("groupName")
	password := ctx.FormValue("password")
	public := ctx.FormValue("public") != ""

	err := groupOperations.CreateGroup(groupName, userID, password, public)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Group created.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//JoinPublicGroup - create a new group with groupName and user id as the owner
func JoinPublicGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	password := ctx.FormValue("password")

	err := groupOperations.JoinPublicGroup(userID, groupID, password)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Group joined.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// SearchPublicGroups -
func SearchPublicGroups(ctx echo.Context) error {
	groupName := ctx.FormValue("groupName")

	groups, err := groupOperations.SearchPublicGroups(groupName)

	switch err {
	case nil:
		return ctx.JSON(200, groups)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// GetUserGroups -
func GetUserGroups(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)

	groups, err := groupOperations.GetUserGroups(userID)

	switch err {
	case nil:
		return ctx.JSON(200, groups)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// GetGroupUsers -
func GetGroupUsers(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	userList, err := groupOperations.GetGroupUsers(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, userList)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//StoreMessage - store a message in a group
func StoreMessage(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	message := ctx.FormValue("message")

	err := groupOperations.StoreUserGroupMessages(groupID, userID, message)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Message stored.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//GetMessages - get all messages for group
func GetMessages(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	messages, err := groupOperations.GetUserGroupMessages(groupID, userID)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"messages": messages})
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
