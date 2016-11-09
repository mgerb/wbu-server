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
func CreateGroup(groupName string, userID string, usersName string) error {

	//DO VALIDATION
	if !regexp.MustCompile(regex.GROUP_NAME).MatchString(groupName) {
		return errors.New("Invalid group name.")
	}

	groupExists := db.Client.HExists(groupModel.GROUP_ID(), groupName).Val()

	if !groupExists {
		return errors.New("Group already exists.")
	}

	temp, _ := db.Client.Incr(groupModel.GROUP_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.HSet(groupModel.GROUP_ID(), groupName, newID)

	//store group hash
	pipe.HMSet(groupModel.GROUP_HASH(newID), groupModel.GROUP_HASH_MAP(groupName, userID))

	pipe.HSet(groupModel.GROUP_MEMBERS(newID), userID, usersName)
	pipe.HSet(userModel.USER_GROUPS(userID), newID, groupName)
	pipe.HIncrBy(userModel.USER_HASH(userID), "adminGroupCount", 1)

	_, returnError := pipe.Exec()

	return returnError
}

//GetGroupMembers - returns string array of group members - userID/userName
func GetMembers(userID string, groupID string) (map[string]string, error) {

	userIsMember := db.Client.HExists(groupModel.GROUP_MEMBERS(groupID), userID).Val()

	if !userIsMember {
		return map[string]string{}, errors.New("You are not a member of this group")
	}

	return db.Client.HGetAll(groupModel.GROUP_MEMBERS(groupID)).Result()
}

func InviteToGroup(groupOwnerID string, groupID string, groupName string, invUserID string) error {
	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//check if inviter is group owner
	tempOwnerID := pipe.HGet(groupModel.GROUP_HASH(groupID), "owner")

	//check if user exists
	tempUserExists := pipe.Exists(userModel.USER_HASH(invUserID))

	//check if user already exists in group
	tempUserExistsInGroup := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), invUserID)

	//check if user already has an invite to the group that is pending
	tempUserHasPendingInvite := pipe.HExists(userModel.USER_GROUP_INVITES(invUserID), groupID)

	_, err := pipe.Exec()

	if err != nil {
		return errors.New("Error inviting user.")
	}

	ownerID := tempOwnerID.Val()

	if ownerID != groupOwnerID {
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

	addInviteErr := db.Client.HSet(userModel.USER_GROUP_INVITES(invUserID), groupID, groupName).Err()

	if addInviteErr != nil {
		return addInviteErr
	}

	return nil
}

//TODO-----------------------------------------------------------------
//TODO-----------------------------------------------------------------
