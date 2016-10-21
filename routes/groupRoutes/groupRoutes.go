package groupRoutes

import (
	"../../operations/groupOperations"
	"github.com/kataras/iris"
)

//CreateGroup - create a new group with groupname and user id as the owner
func CreateGroup(ctx *iris.Context) {
	userid := ctx.Get("id")
	username := ctx.Get("username")
	groupname := ctx.PostValue("groupname")
	
	err := groupOperations.CreateGroup(groupname, userid.(string), username.(string))
	
	if err == nil {
		ctx.JSON(200, `{"message": "Group created"}`)
	} else {
		ctx.JSON(500, `{"message": "` + err.Error() + `"}`)
	}
}
