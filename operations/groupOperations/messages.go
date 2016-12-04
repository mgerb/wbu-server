package groupOperations

import (
	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"errors"
	redis "gopkg.in/redis.v5"
	"time"
)

func StoreUserGroupMessages(groupID string, userID string, message string) error {

	//do validation before running redis script
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("Invalid message length.")
	}

	luaScript := `
		local groupIDKey = KEYS[1]
		local userIDKey = KEYS[2]
		local userGrpMsgKey = ARGV[1]
		local userID = ARGV[2]
		local groupID = ARGV[3]
		local timeStamp = ARGV[4]
		local message = ARGV[5]
		local oneMonth = 2592000
		
		
		--check if user exists in group
		if not redis.call("HGET", groupIDKey, userID) then
			return redis.error_reply("You are not in this group")	
		end
		
		local fullName = redis.call("HGET", userIDKey, "fullName")
		
		-- get user full name
		if not fullName then
			return redis.error_reply("User does not exist")	
		end
		
		local fullMessage = userID .. "/" .. fullName .. "/" .. timeStamp .. "/" .. message
		
		-- get all the members in the group
		local members = redis.call("HGETALL", groupIDKey)
		
		-- cycle through each key in the group member hash
		for i = 1, #members, 2 do
			redis.call("SADD", userGrpMsgKey  .. members[i] .. ":" .. groupID, fullMessage)
			-- reset the expire time for each message
			redis.call("EXPIRE", userGrpMsgKey  .. members[i] .. ":" .. groupID, oneMonth)
		end
		
		return "Success"
	`

	script := redis.NewScript(luaScript)

	return script.Run(db.Client, []string{groupModel.GROUP_MEMBERS(groupID), userModel.USER_HASH(userID)},
		userModel.USER_GROUP_MESSAGE_KEY(),
		userID,
		groupID,
		time.Now().Unix(),
		message).Err()
}

func GetUserGroupMessages(groupID string, userID string) (interface{}, error) {
	luaScript := `
		local key = KEYS[1]
		
		local messages = redis.call("SMEMBERS", key)
		
		-- remove the messages after retreiving
		redis.call("DEL", key)
		
		return messages
	`

	script := redis.NewScript(luaScript)

	return script.Run(db.Client, []string{userModel.USER_GROUP_MESSAGES(userID, groupID)}, 0).Result()
}

/*
DEPRECATED MESSAGE FUNCTIONS

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
*/
