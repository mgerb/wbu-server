package groupOperations

import (
	"../../db"
	"../../model/groupModel"
	"errors"
	"strconv"
	"time"
)

//set max number of messages stored at any one point in time
var maxMessages int64 = 99

//StoreMessage - store a message and a group list - maximum of 100 messages stored at any point
func StoreMessage(groupID string, userID string, userName string, message string) error {

	//DO VALIDATION
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("Invalid message length.")
	}
	//check if user exists in group before storing message
	if !UserIsMember(userID, userName, groupID) {
		return errors.New("User is not in group.")
	}

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.LPush(groupModel.GROUP_MESSAGE(groupID), userID+"/"+userName+"/"+strconv.FormatInt(time.Now().Unix(), 10)+"/"+message)
	pipe.LTrim(groupModel.GROUP_MESSAGE(groupID), 0, maxMessages)

	_, err := pipe.Exec()

	return err
}

func GetMessages(userID string, userName string, groupID string) ([]string, error) {

	//DO VALIDATION
	//check if user exists in group before storing message
	if !UserIsMember(userID, userName, groupID) {
		return []string{}, errors.New("User is not in group.")
	}

	return db.Client.LRange(groupModel.GROUP_MESSAGE(groupID), 0, -1).Result()

	/*
		if err == nil {
			json, jsonError := json.Marshal(messages)

			if jsonError == nil {
				return string(json), nil
			} else {
				return "", errors.New("json error")
			}
		} else {
			return "", errors.New("Error retrieving messages.")
		}
	*/
}
