package groupOperations

import (
	"errors"
	"strconv"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
)

//CreateGroup - store username/password in hash
func CreateGroup(groupname string, userID string, username string) error {
	_, err := GetGroupID(groupname)
	if err != nil {
		temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
		newID := strconv.FormatInt(temp, 10)

		pipe := db.Client.Pipeline()
		defer pipe.Close()

		pipe.Set(groupModel.GROUP_ID(groupname), newID, 0)
		pipe.HMSet(groupModel.GROUP_HASH(newID), map[string]string{
			"groupname": groupname,
			"owner":     userID,
		})

		pipe.SAdd(groupModel.GROUP_MEMBERS(newID), userID+"/"+username)
		pipe.SAdd(userModel.USER_GROUPS(userID), newID+"/"+groupname)

		_, err = pipe.Exec()

		return err
	}

	return errors.New("group already exists")
}

//GetGroupID - get the group id - check if group exists
func GetGroupID(groupname string) (string, error) {
	return db.Client.Get(groupModel.GROUP_ID(groupname)).Result()
}

//GetGroupMembers - returns string array of group members - userID/userName
func GetGroupMembers(groupID string) []string {
	result, _ := db.Client.SMembers(groupModel.GROUP_MEMBERS(groupID)).Result()
	return result
}

//UserIsMember - returns nil if user is member
func UserIsMember(userID string, groupID string) error {
	_, _, err := db.Client.SScan(groupModel.GROUP_MEMBERS(groupID), 0, userID+"/*", 1).Result()
	return err
}

func InviteToGroup(groupOwnerId string, groupid string, inviteduserID string) error {
	return errors.New("TODO")
}

//TODO-----------------------------------------------------------------
func JoinGroup(userID string) error {
	return errors.New("TODO")
}

func LeaveGroup(userID string, groupid string) error {
	return errors.New("TODO")
}

//TODO-----------------------------------------------------------------
