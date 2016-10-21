package groupRoutes

import (
	"../../operations/groupOperations"
	"../../utils/response"
	"github.com/kataras/iris"
)

//CreateGroup - create a new group with groupName and user id as the owner
func CreateGroup(ctx *iris.Context) {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)
	groupName := ctx.PostValue("groupName")

	err := groupOperations.CreateGroup(groupName, userID, userName)

	if err == nil {
		ctx.JSON(200, response.Json("Group created.", response.SUCCESS))
	} else {
		ctx.JSON(500, `{"message": "`+err.Error()+`"}`)
	}
}

//StoreMessage - store a message in a group
func StoreMessage(ctx *iris.Context) {
	userID := ctx.Get("userID").(string)
	userName := ctx.Get("userName").(string)

	groupID := ctx.PostValue("groupID")
	message := ctx.PostValue("message")

	err := groupOperations.StoreMessage(groupID, userID, userName, message)

	if err == nil {
		ctx.JSON(200, response.Json("Message received.", response.SUCCESS))
	} else {
		ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
