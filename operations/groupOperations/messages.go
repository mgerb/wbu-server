package groupOperations

import (
	"errors"
	"strconv"
	"time"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
)

//set max number of messages stored at any one point in time
var maxMessages int64 = 99

//StoreMessage - store a message and a group list - maximum of 100 messages stored at any point
func StoreMessage(groupID string, userID string, message string) error {

	//TODO - PUSH NOTIFICATIONS
	//-------------------------
	//-------------------------
	//-------------------------
	//-------------------------
	//-------------------------
	//-------------------------

	//DO VALIDATION
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("Invalid message length.")
	}

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	tempFullName := pipe.HGet(userModel.USER_HASH(userID), "fullName")
	userIsMember := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), userID)

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return errors.New("Pipe error")
	}

	if !userIsMember.Val() {
		return errors.New("You are not a member of this group")
	}

	fullName := tempFullName.Val()

	pipe.LPush(groupModel.GROUP_MESSAGES(groupID), userID+"/"+fullName+"/"+strconv.FormatInt(time.Now().Unix(), 10)+"/"+message)
	pipe.LTrim(groupModel.GROUP_MESSAGES(groupID), 0, maxMessages)

	_, err := pipe.Exec()

	if err != nil {
		return errors.New("Error storing message.")
	}

	return err
}

func GetMessages(userID string, groupID string) ([]string, error) {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	groupExists := pipe.Exists(groupModel.GROUP_HASH(groupID))
	userHasGroup := pipe.HExists(userModel.USER_GROUPS(userID), groupID)

	_, err := pipe.Exec()

	if err != nil {
		return []string{}, errors.New("Database error")
	}

	if !userHasGroup.Val() {
		return []string{}, errors.New("You are not a member of this group")
	} else if !groupExists.Val() && userHasGroup.Val() {
		//delete group if group does not exist anymore, but user still has group in user group hash
		db.Client.HDel(userModel.USER_GROUPS(userID), groupID)
		return []string{}, errors.New("Group has been removed.")
	}

	return db.Client.LRange(groupModel.GROUP_MESSAGES(groupID), 0, -1).Result()
}
