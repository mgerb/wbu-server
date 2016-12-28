package groupOperations

import (
	"errors"
	"regexp"
	"strconv"

	"../lua"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils/regex"
	redis "gopkg.in/redis.v5"
)

//CreateGroup - store userName/password in hash
func CreateGroup(groupName string, userID string) error {

	//DO VALIDATION
	if !regexp.MustCompile(regex.GROUP_NAME).MatchString(groupName) {
		return errors.New("Invalid group name.")
	}

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	tempFullName := pipe.HGet(userModel.USER_HASH(userID), "fullName")
	groupExists := pipe.HExists(groupModel.GROUP_ID(), groupName)

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return errors.New("Pipe error.")
	}

	if groupExists.Val() {
		return errors.New("Group already exists.")
	}

	fullName := tempFullName.Val()

	temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe.HSet(groupModel.GROUP_ID(), groupName, newID)

	//store group hash
	pipe.HMSet(groupModel.GROUP_HASH(newID), groupModel.GROUP_HASH_MAP(groupName, userID))

	pipe.HSet(groupModel.GROUP_MEMBERS(newID), userID, fullName)
	pipe.HSet(userModel.USER_GROUPS(userID), newID, groupName)
	pipe.HIncrBy(userModel.USER_HASH(userID), "adminGroupCount", 1)

	_, returnError := pipe.Exec()

	return returnError
}

//GetGroupMembers - returns string array of group members - userID/userName
func GetGroupMembers(userID string, groupID string) (interface{}, error) {

	script := redis.NewScript(lua.Use("GetGroupMembers.lua"))

	return script.Run(db.Client, []string{
		groupModel.GROUP_HASH(groupID),
		groupModel.GROUP_MEMBERS(groupID),
	},
		userID,
	).Result()
}

func InviteToGroup(groupOwnerID string, groupID string, invUserID string) error {
	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//get group information
	groupInfo := pipe.HGetAll(groupModel.GROUP_HASH(groupID))

	//check if user exists
	tempUserExists := pipe.Exists(userModel.USER_HASH(invUserID))

	//check if user already exists in group
	tempUserExistsInGroup := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), invUserID)

	//check if user already has an invite to the group that is pending
	tempUserHasPendingInvite := pipe.HExists(userModel.USER_GROUP_INVITES(invUserID), groupID)

	_, err := pipe.Exec()

	storedGroupOwnerID := groupInfo.Val()["owner"]
	storedGroupName := groupInfo.Val()["groupName"]

	if err != nil {
		return errors.New("Error inviting user.")
	}

	if groupOwnerID != storedGroupOwnerID {
		return errors.New("User does not have permission.")
	}

	userExists := tempUserExists.Val()

	if !userExists {
		return errors.New("User does not exist.")
	}

	userHasPendingInvite := tempUserHasPendingInvite.Val()

	if userHasPendingInvite {
		return errors.New("User has pending invite.")
	}

	userExistsInGroup := tempUserExistsInGroup.Val()

	if userExistsInGroup {
		return errors.New("User is already in group.")
	}

	addInviteErr := db.Client.HSet(userModel.USER_GROUP_INVITES(invUserID), groupID, storedGroupName).Err()

	if addInviteErr != nil {
		return addInviteErr
	}

	return nil
}

func JoinGroup(userID string, groupID string) error {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	userInfo := pipe.HGetAll(userModel.USER_HASH(userID))
	groupInfo := pipe.HGetAll(groupModel.GROUP_HASH(groupID))

	groupExists := pipe.Exists(groupModel.GROUP_HASH(groupID))
	userHasInvite := pipe.HExists(userModel.USER_GROUP_INVITES(userID), groupID)

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return err_pipe1
	}

	if !groupExists.Val() {
		//delete user invite if group has been deleted and they still have an invite
		if userHasInvite.Val() {
			db.Client.HDel(userModel.USER_GROUP_INVITES(userID), groupID)
		}
		return errors.New("Group does not exist.")
	}

	if !userHasInvite.Val() {
		return errors.New("User does not have an invite.")
	}

	fullName := userInfo.Val()["fullName"]
	groupName := groupInfo.Val()["groupName"]

	pipe.HSet(groupModel.GROUP_MEMBERS(groupID), userID, fullName)
	pipe.HSet(userModel.USER_GROUPS(userID), groupID, groupName)
	pipe.HDel(userModel.USER_GROUP_INVITES(userID), groupID)

	_, err_pipe2 := pipe.Exec()

	if err_pipe2 != nil {
		return err_pipe2
	}

	return nil
}

func LeaveGroup(userID string, groupID string) error {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	userExistsInGroup := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), userID)
	groupExistsInUser := pipe.HExists(userModel.USER_GROUPS(userID), groupID)

	_, err1 := pipe.Exec()

	if err1 != nil {
		return errors.New("Pipe error.")
	}

	if !userExistsInGroup.Val() && !groupExistsInUser.Val() {
		return errors.New("You are not a member of this group.")
	}

	pipe.HDel(groupModel.GROUP_MEMBERS(groupID), userID)
	pipe.HDel(userModel.USER_GROUPS(userID), groupID)

	_, err2 := pipe.Exec()

	return err2
}

func DeleteGroup(userID string, groupID string) error {

	script := redis.NewScript(lua.Use("DeleteGroup.lua"))

	return script.Run(db.Client, []string{
		userModel.USER_HASH(userID),
		groupModel.GROUP_ID(),
		groupModel.GROUP_HASH(groupID),
		groupModel.GROUP_MEMBERS(groupID),
		groupModel.GROUP_MESSAGES(groupID),
		groupModel.GROUP_GEO(groupID),
		groupModel.GROUP_LOCATIONS(groupID),
	},
		userID,
		groupID,
		userModel.USER_GROUPS_KEY(),
		userModel.USER_GROUP_MESSAGES_KEY(),
	).Err()
}
