package groupRoutes

import (
	"log"

	"../../operations/groupOperations"
	"github.com/kataras/iris"
)

//create a new group with groupname and user id as the owner
func CreateGroup(ctx *iris.Context) {
	userid := ctx.Get("id")
	groupname := ctx.PostValue("groupname")

	if groupname != "" {
		err := groupOperations.CreateGroup(groupname, userid.(string))
		if err != nil {
			log.Println(err)
			ctx.JSON(500, `{"message": "Group already exists"}`)
		} else {
			ctx.JSON(200, `{"message": "Group created"}`)
		}
	} else {
		ctx.JSON(500, `{"message": "Invalid group name."}`)
	}
}
