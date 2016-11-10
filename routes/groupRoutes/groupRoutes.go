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

	err := groupOperations.CreateGroup(groupName, userID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Group created.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetMembers(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	members, err := groupOperations.GetMembers(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, members)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//StoreMessage - store a message in a group
func StoreMessage(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	message := ctx.FormValue("message")

	err := groupOperations.StoreMessage(groupID, userID, message)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Message received.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

//GetMessages - get all messages for group
func GetMessages(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	messages, err := groupOperations.GetMessages(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"messages": messages})
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func InviteUser(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	invUserID := ctx.FormValue("invUserID")

	err := groupOperations.InviteToGroup(userID, groupID, invUserID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("User invited.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func JoinGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")

	err := groupOperations.JoinGroup(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Joined group.", response.SUCCESS))

	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
