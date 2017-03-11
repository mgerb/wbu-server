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

func GetGroupMembers(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	members, err := groupOperations.GetGroupMembers(userID, groupID)

	switch err {
	case nil:
		//maybe change this end point in future
		//return ctx.JSON(200, members)
		return ctx.Blob(200, response.JSON_HEADER, []byte(members.(string)))
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

func LeaveGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")

	err := groupOperations.LeaveGroup(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Left group.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func DeleteGroup(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")

	err := groupOperations.DeleteGroup(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Group deleted.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func StoreGeoLocation(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	latitude := ctx.FormValue("latitude")
	longitude := ctx.FormValue("longitude")

	err := groupOperations.StoreGeoLocation(userID, groupID, latitude, longitude)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Location stored.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

func GetGeoLocations(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	geoLocations, err := groupOperations.GetGeoLocations(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, map[string]interface{}{"geoLocations": geoLocations})
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
