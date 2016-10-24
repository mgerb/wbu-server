package groupOperations

import (
	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils/regex"
	"errors"
	"regexp"
	"strconv"
)

//CreateGroup - store userName/password in hash
func CreateGroup(groupName string, userID string, userName string) error {

	//DO VALIDATION
	if !regexp.MustCompile(regex.GROUP_NAME).MatchString(groupName) {
		return errors.New("Invalid group name.")
	}

	_, err := GetGroupID(groupName)
	if err == nil {
		return errors.New("Group already exists.")
	}

	temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.Set(groupModel.GROUP_ID(groupName), newID, 0)

	//store group hash
	pipe.HMSet(groupModel.GROUP_HASH(newID), map[string]string{
		"groupName": groupName,
		"owner":     userID,
	})

	pipe.SAdd(groupModel.GROUP_MEMBERS(newID), userID+"/"+userName)
	pipe.SAdd(userModel.USER_GROUPS(userID), newID+"/"+groupName)

	_, returnError := pipe.Exec()

	return returnError
}

//GetGroupID - get the group id - check if group exists
func GetGroupID(groupName string) (string, error) {
	return db.Client.Get(groupModel.GROUP_ID(groupName)).Result()
}

//GetGroupMembers - returns string array of group members - userID/userName
func GetMembers(userID string, userName string, groupID string) ([]string, error) {

	if !UserIsMember(userID, userName, groupID) {
		return []string{}, errors.New("You are not a member of this group.")
	}

	return db.Client.SMembers(groupModel.GROUP_MEMBERS(groupID)).Result()
}

//UserIsMember - returns true if user is member
func UserIsMember(userID string, userName string, groupID string) bool {
	isMember, err := db.Client.SIsMember(groupModel.GROUP_MEMBERS(groupID), userID+"/"+userName).Result()

	if err == nil {
		return isMember
	} else {
		return false
	}
}

//TODO-----------------------------------------------------------------
func InviteToGroup(groupOwnerId string, groupID string, inviteduserID string) error {
	return errors.New("TODO")
}

//TODO-----------------------------------------------------------------
