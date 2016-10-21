package groupOperations

import (
	"errors"
	"regexp"
	"strconv"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils"
)

//CreateGroup - store username/password in hash
func CreateGroup(groupname string, userID string, username string) error {
	
	//DO VALIDATION
	if !regexp.MustCompile(utils.GroupnameRegex).MatchString(username){
		return errors.New("Invalid group name.")
	}
	
	_, err := GetGroupID(groupname)
	if err == nil {
		return errors.New("Group already exists.")
	}

	temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.Set(groupModel.GROUP_ID(groupname), newID, 0)

	//store group hash
	pipe.HMSet(groupModel.GROUP_HASH(newID), map[string]string{
		"groupname": groupname,
		"owner":     userID,
	})

	pipe.SAdd(groupModel.GROUP_MEMBERS(newID), userID+"/"+username)
	pipe.SAdd(userModel.USER_GROUPS(userID), newID+"/"+groupname)

	_, returnError := pipe.Exec()

	return returnError
}

//GetGroupID - get the group id - check if group exists
func GetGroupID(groupname string) (string, error) {
	return db.Client.Get(groupModel.GROUP_ID(groupname)).Result()
}

//GetGroupMembers - returns string array of group members - userID/userName
func GetGroupMembers(groupID string) ([]string, error) {
	return db.Client.SMembers(groupModel.GROUP_MEMBERS(groupID)).Result()
}

//UserIsMember - returns true if user is member
func UserIsMember(userID string, groupID string) error {
	_, _, err := db.Client.SScan(groupModel.GROUP_MEMBERS(groupID), 0, userID+"/*", 1).Result()
	return err
}

//TODO-----------------------------------------------------------------
func InviteToGroup(groupOwnerId string, groupid string, inviteduserID string) error {
	return errors.New("TODO")
}

//TODO-----------------------------------------------------------------
