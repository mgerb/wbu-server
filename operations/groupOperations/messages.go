package groupOperations

import (
	"errors"
	"time"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../lua"
	redis "gopkg.in/redis.v5"
)

//StoreUserGroupMessages - store a users messages in a group
func StoreUserGroupMessages(groupID string, userID string, message string) error {

	//do validation before running redis script
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("Invalid message length.")
	}

	script := redis.NewScript(lua.Use("StoreUserGroupMessages.lua"))

	return script.Run(db.Client, []string{groupModel.GROUP_MEMBERS(groupID), userModel.USER_HASH(userID)},
		userModel.USER_GROUP_MESSAGES_KEY(),
		userID,
		groupID,
		time.Now().Unix(),
		message).Err()
}

//GetUserGroupMessages - return user messages for a group
func GetUserGroupMessages(groupID string, userID string) (interface{}, error) {

	script := redis.NewScript(lua.Use("GetUserGroupMessages.lua"))

	return script.Run(db.Client, []string{userModel.USER_GROUP_MESSAGES(userID, groupID)}, 0).Result()
}
