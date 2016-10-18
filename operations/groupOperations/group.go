package groupOperations

import (
	"errors"
	"strconv"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
)

//CreateGroup - store username/password in hash
func CreateGroup(groupname string, userID string) error {
	if Exists(groupname) != true {
		temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
		newID := strconv.FormatInt(temp, 10)

		db.Client.Set(groupModel.GROUP_NAME(groupname), newID, 0)
		db.Client.HMSet(groupModel.GROUP_HASH(newID), map[string]string{
			"groupname": groupname,
			"owner":     userID,
		})

		db.Client.SAdd(groupModel.GROUP_MEMBERS(newID), userID)
		db.Client.SAdd(groupModel.GROUP_MEMBERS(newID), 123345)
		db.Client.SAdd(userModel.USER_GROUPS(userID), newID)
		return nil
	}
	return errors.New("group already exists")
}

//Exists - check if group exists in redis - return boolean
func Exists(groupname string) bool {
	return db.Client.Get(groupModel.GROUP_NAME(groupname)).Err() == nil
}

//TODO-----------------------------------------------------------------
func JoinGroup(userID string) error {
	return errors.New("TODO")
}

func LeaveGroup(userID string, groupid string) error {
	return errors.New("TODO")
}

func InviteToGroup(groupOwnerId string, groupid string, inviteduserID string) error {
	return errors.New("TODO")
}

//TODO-----------------------------------------------------------------
