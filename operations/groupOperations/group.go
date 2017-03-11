package groupOperations

import (
	"errors"
	"log"
	"regexp"

	"../lua"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils"
	"../../utils/regex"
	redis "gopkg.in/redis.v5"
)

//CreateGroup - store userName/password in hash
func CreateGroup(groupName string, userID string, password string, public bool) error {

	// validate group name
	if !regexp.MustCompile(regex.GROUP_NAME).MatchString(groupName) {
		return errors.New("Invalid group name.")
	}

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	defer tx.Commit()

	//check if the group already exists
	var groupExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE "name" = ? AND "ownerID" = ?);`, groupName, userID).Scan(&groupExists)

	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	} else if groupExists {
		return errors.New("Group already exists.")
	}

	if public {

		// if users sets a password for the group
		if password != "" {
			// validate password
			if len(password) < 5 {
				return errors.New("Password must be more than 5 characters.")
			}

			// hash the password before storing in the database
			passwordHash, err := utils.GenerateHash(password)

			if err != nil {
				log.Println(err)
				return errors.New("Error hashing password")
			}

			_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID, userCount, public, password, locked) VALUES (?, ?, ?, ?, ?, ?);`, groupName, userID, 1, 1, passwordHash, 1)
		} else {
			_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID, userCount, public) VALUES (?, ?, ?, ?);`, groupName, userID, 1, 1)
		}
	} else {
		_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID, userCount) VALUES (?, ?, ?);`, groupName, userID, 1)
	}

	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	return nil
}

func SearchPublicGroups(groupName string) ([]*groupModel.Group, error) {
	groupList := []*groupModel.Group{}

	// get user information
	rows, err := db.SQL.Query(`
		SELECT g.id, g.name, u.email, g.userCount, g.locked, u.firstName, u.lastName
		FROM "Group" AS g INNER JOIN "USER" AS u ON g.ownerID = u.id
		WHERE g.public = 1 AND g.name LIKE ?;
		`,
		"%"+groupName+"%")

	if err != nil {
		log.Println(err)
		return []*groupModel.Group{}, errors.New("Database error.")
	}

	defer rows.Close()

	for rows.Next() {
		newGroup := &groupModel.Group{}
		var firstName string
		var lastName string
		err := rows.Scan(&newGroup.ID, &newGroup.Name, &newGroup.OwnerEmail, &newGroup.UserCount, &newGroup.Locked, &firstName, &lastName)
		newGroup.OwnerName = firstName + " " + lastName

		if err != nil {
			log.Println(err)
			return []*groupModel.Group{}, errors.New("Database error.")
		}

		groupList = append(groupList, newGroup)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*groupModel.Group{}, errors.New("Row error.")
	}

	return groupList, nil
}

// get all the group for a specific user
func GetUserGroups(userID string) ([]*groupModel.Group, error) {
	groupList := []*groupModel.Group{}

	// get user information
	rows, err := db.SQL.Query(`
		SELECT g.id, g.name, g.ownerID, g.userCount, g.locked FROM "Group" AS g
		INNER JOIN "UserGroup" AS ug ON g.id = ug.groupID
		INNER JOIN "User" AS u ON ug.userID = u.id
		WHERE u.id = ?;
		`,
		"%"+userID+"%")

	if err != nil {
		log.Println(err)
		return []*groupModel.Group{}, errors.New("Database error.")
	}

	defer rows.Close()

	for rows.Next() {
		newGroup := &groupModel.Group{}
		err := rows.Scan(&newGroup.ID, &newGroup.Name, &newGroup.OwnerID, &newGroup.UserCount, &newGroup.Locked)

		if err != nil {
			log.Println(err)
			return []*groupModel.Group{}, errors.New("Database error.")
		}

		groupList = append(groupList, newGroup)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*groupModel.Group{}, errors.New("Row error.")
	}

	return groupList, nil
}

// TODO
func JoinGroupWithPassword(userID string, ownerID string, groupID int, password string) {
	/*
		tx, err := db.SQL.Begin()
		if err != nil {
			log.Println(err)
			return errors.New("Database error.")
		}

		defer tx.Commit()

		//check if the group already exists
		var groupExists bool
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE "name" = ? AND "ownerID" = ?);`, groupName, userID).Scan(&groupExists)

		if err != nil {
			log.Println(err)
			return errors.New("Database error.")
		} else if groupExists {
			return errors.New("Group already exists.")
		}
	*/

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
