package groupRoutes

import (
	"../../operations/groupOperations"
	"../../utils/response"
	"github.com/labstack/echo"
)

//CreateGroup - create a new group with groupName and user id as the owner
func CreateGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)
	groupName := ctx.FormValue("groupName")

	err := groupOperations.CreateGroup(groupName, userID, userName)

	if err == nil {
		return ctx.JSON(200, response.Json("Group created.", response.SUCCESS))
	} else {
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetMembers(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)
	groupID := ctx.Param("groupID")

	members, err := groupOperations.GetMembers(userID, userName, groupID)

	if err == nil {
		return ctx.JSON(200, map[string]interface{}{"members": members})
	} else {
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//StoreMessage - store a message in a group
func StoreMessage(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)

	groupID := ctx.Param("groupID")
	message := ctx.FormValue("message")

	err := groupOperations.StoreMessage(groupID, userID, userName, message)

	if err == nil {
		return ctx.JSON(200, response.Json("Message received.", response.SUCCESS))
	} else {
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//GetMessages - get all messages for group
func GetMessages(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)
	groupID := ctx.Param("groupID")

	messages, err := groupOperations.GetMessages(userID, userName, groupID)

	if err == nil {
		return ctx.JSON(200, map[string]interface{}{"messages": messages})
	} else {
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func InviteToGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	groupName := ctx.FormValue("groupName")
	invUserID := ctx.FormValue("invUserID")
	invUserName := ctx.FormValue("invUserName")

	err := groupOperations.InviteToGroup(userID, groupID, groupName, invUserID, invUserName)

	if err == nil {
		return ctx.JSON(200, map[string]interface{}{"message": "success"})
	} else {
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
