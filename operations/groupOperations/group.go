package groupOperations

import (
	"database/sql"
	"errors"
	"log"
	"regexp"

	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/model"
	"github.com/mgerb/wbu-server/utils"
	"github.com/mgerb/wbu-server/utils/fcm"
	"github.com/mgerb/wbu-server/utils/regex"
	"golang.org/x/crypto/bcrypt"
)

//CreateGroup - create new group in Group table - also add owner to UserGroup table
func CreateGroup(name string, userID string, password string, public bool) error {

	// validate group name
	if !regexp.MustCompile(regex.GROUP_NAME).MatchString(name) {
		return errors.New("invalid group name")
	}

	// start the SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// commit transaction when function returns
	defer tx.Commit()

	//check if the group already exists
	var groupExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE "name" = ? AND "ownerID" = ?);`, name, userID).Scan(&groupExists)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	if groupExists {
		return errors.New("group already exists")
	}

	if public {
		// if users sets a password for the group
		if password != "" {

			// hash the password before storing in the database
			passwordHash, err := utils.GenerateHash(password)

			if err != nil {
				log.Println(err)
				return errors.New("Error hashing password")
			}

			_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID, public, password) VALUES (?, ?, ?, ?);`, name, userID, 1, passwordHash)
		} else {
			_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID, public) VALUES (?, ?, ?);`, name, userID, 1)
		}
	} else { // if private group
		_, err = tx.Exec(`INSERT INTO "Group" (name, ownerID) VALUES (?, ?);`, name, userID)
	}

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	var groupID string
	err = tx.QueryRow(`SELECT id FROM "GROUP" WHERE name = ? AND ownerID = ?;`, name, userID).Scan(&groupID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// insert id's into UserGroup table
	err = insertUserGroup(tx, userID, groupID)

	if err != nil {
		log.Println(err)
		// roll back the transaction if we get an error inserting into UserGroup
		tx.Rollback()
		return errors.New("database error")
	}

	return nil
}

// SearchPublicGroups - matches group by group name
func SearchPublicGroups(searchTerm string) ([]*model.Group, error) {
	groupList := []*model.Group{}

	// query groups - join with UserGroup and User tables to get the group owner information
	rows, err := db.SQL.Query(`
		SELECT g.id, g.name, u.email, g.userCount, g.password, u.firstName, u.lastName
		FROM "Group" AS g INNER JOIN "USER" AS u ON g.ownerID = u.id
		WHERE g.public = 1 AND (g.name LIKE ? OR u.email LIKE lower(?)) LIMIT 20;`,
		"%"+searchTerm+"%", "%"+searchTerm+"%")

	if err != nil {
		log.Println(err)
		return []*model.Group{}, errors.New("database error")
	}

	defer rows.Close()

	// map query to object list
	for rows.Next() {
		newGroup := &model.Group{}
		var firstName string
		var lastName string
		err := rows.Scan(&newGroup.ID, &newGroup.Name, &newGroup.OwnerEmail, &newGroup.UserCount, &newGroup.Password, &firstName, &lastName)
		newGroup.OwnerName = firstName + " " + lastName

		// group is locked if password is not null
		newGroup.Locked = newGroup.Password.Valid

		if err != nil {
			log.Println(err)
			return []*model.Group{}, errors.New("database error")
		}

		groupList = append(groupList, newGroup)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.Group{}, errors.New("row error")
	}

	return groupList, nil
}

//GetUserGroups - get group list for a specific user
func GetUserGroups(userID string) ([]*model.Group, error) {
	groupList := []*model.Group{}

	// get group list - join Group on UserGroup and User
	rows, err := db.SQL.Query(`
		SELECT g.id, g.name, g.ownerID, u.firstName, u.lastName, u.email, g.userCount, g.password , g.public FROM "Group" AS g
		INNER JOIN "UserGroup" AS ug ON g.id = ug.groupID
		INNER JOIN "User" AS u ON ug.userID = u.id
		WHERE u.id = ?;`, userID)

	if err != nil {
		log.Println(err)
		return []*model.Group{}, errors.New("database error")
	}

	defer rows.Close()

	// map query to group object list
	for rows.Next() {
		newGroup := &model.Group{}
		var firstName, lastName string
		err := rows.Scan(&newGroup.ID, &newGroup.Name, &newGroup.OwnerID, &firstName, &lastName, &newGroup.OwnerEmail, &newGroup.UserCount, &newGroup.Password, &newGroup.Public)
		newGroup.OwnerName = firstName + " " + lastName

		// group is locked if password is not null
		newGroup.Locked = newGroup.Password.Valid

		if err != nil {
			log.Println(err)
			return []*model.Group{}, errors.New("database error")
		}

		groupList = append(groupList, newGroup)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.Group{}, errors.New("row error")
	}

	return groupList, nil
}

//GetGroupUsers - list users for single group
func GetGroupUsers(userID string, groupID string) ([]*model.User, error) {

	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("database error")
	}

	defer tx.Commit()

	// check if user exists in group before getting group users
	// check ID's in UserGroup table
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE "groupID" = ? AND "userID" = ?);`, groupID, userID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("database error")
	} else if !userExistsInGroup {
		return []*model.User{}, errors.New("user not in group")
	}

	// get list of users if user exists in group
	userList := []*model.User{}

	// get user information
	rows, err := tx.Query(`
		SELECT u.id, u.email, u.firstName, u.lastName FROM "User" AS u
		INNER JOIN "UserGroup" AS ug ON u.id = ug.userID
		INNER JOIN "Group" AS g ON ug.groupID = g.id
		WHERE g.id = ?;`, groupID)

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("database error")
	}

	defer rows.Close()

	for rows.Next() {
		newUser := &model.User{}
		err := rows.Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName)

		if err != nil {
			log.Println(err)
			return []*model.User{}, errors.New("database error")
		}

		userList = append(userList, newUser)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("row error")
	}

	return userList, nil
}

// JoinPublicGroup -
func JoinPublicGroup(userID string, groupID string, password string) error {

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	newGroup := &model.Group{}

	err = tx.QueryRow(`SELECT password, public FROM "Group" WHERE id = ?;`, groupID).Scan(&newGroup.Password, &newGroup.Public)

	if err != nil {
		// will error if group does not exist
		log.Println(err)
		return errors.New("database error")
	}

	if !newGroup.Public {
		return errors.New("Not a public group.")
	}

	// if group has a password - check if passwords match
	if newGroup.Password.Valid && bcrypt.CompareHashAndPassword([]byte(newGroup.Password.String), []byte(password)) != nil {
		return errors.New("Invalid password.")
	}

	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE "groupID" = ? AND "userID" = ?);`, groupID, userID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	} else if userExistsInGroup {
		return errors.New("Already in group.")
	}

	// insert id's into UserGroup table
	err = insertUserGroup(tx, userID, groupID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("database error")
	}

	// if a user already has a group invite to this group delete it
	_, err = tx.Exec(`DELETE FROM "GroupInvite" WHERE userID = ? AND groupID = ?;`, userID, groupID)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// InviteUserToGroup -
func InviteUserToGroup(ownerID string, userID string, groupID string) error {

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if ownerID = group.ownerID
	// check if user exists
	// check if user already exists in group
	var userIsOwner bool
	var userExists bool
	var userInGroup bool
	var userHasInvite bool
	row := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE "ownerID" = ? AND "id" = ?),
						EXISTS(SELECT 1 FROM "User" WHERE "id" = ?),
						EXISTS(SELECT 1 FROM "UserGroup" WHERE "userID" = ? AND "groupID" = ?),
						EXISTS(SELECT 1 FROM "GroupInvite" WHERE "userID" = ? AND "groupID" = ?);`, ownerID, groupID, userID, userID, groupID, userID, groupID)

	err = row.Scan(&userIsOwner, &userExists, &userInGroup, &userHasInvite)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	if !userIsOwner {
		return errors.New("invalid owner")
	}

	if !userExists {
		return errors.New("invalid user")
	}

	if userInGroup {
		return errors.New("user in group")
	}

	if userHasInvite {
		return errors.New("user has invite")
	}

	// insert (if not exists in table already) into GroupInvite userID and groupID
	_, err = tx.Exec(`INSERT INTO "GroupInvite" (userID, groupID) SELECT ?, ?;`, userID, groupID, userID, groupID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// FCM Notifications
	go fcm.UserInviteNotif(userID)

	return nil
}

// GetGroupInvites -
func GetGroupInvites(userID string) ([]*model.Group, error) {

	rows, err := db.SQL.Query(`SELECT g.ID, g.name, g.userCount, u.firstName, u.lastName, u.email
								FROM "GroupInvite" AS gi INNER JOIN
								"Group" AS g ON gi.groupID = g.id INNER JOIN
								"User" AS u ON g.ownerID = u.id
								WHERE gi.userID = ?;`, userID)

	groupList := []*model.Group{}

	for rows.Next() {
		newGroup := &model.Group{}
		var firstName, lastName string
		err := rows.Scan(&newGroup.ID, &newGroup.Name, &newGroup.UserCount, &firstName, &lastName, &newGroup.OwnerEmail)
		newGroup.OwnerName = firstName + " " + lastName

		if err != nil {
			log.Println(err)
			return []*model.Group{}, errors.New("database error")
		}

		groupList = append(groupList, newGroup)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.Group{}, errors.New("database error")
	}

	return groupList, nil
}

// JoinGroupFromInvite -
func JoinGroupFromInvite(userID string, groupID string) error {

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if user has group invite
	var userHasInvite bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "GroupInvite" WHERE userID = ? AND groupID = ?);`, userID, groupID).Scan(&userHasInvite)

	if !userHasInvite {
		return errors.New("invalid invite")
	} else if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// delete the invite from GroupInvite
	_, err = tx.Exec(`DELETE FROM "GroupInvite" WHERE userID = ? AND groupID = ?;`, userID, groupID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// insert user into group
	err = insertUserGroup(tx, userID, groupID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("database error")
	}

	return nil
}

// DeleteGroupInvite -
func DeleteGroupInvite(userID string, groupID string) error {

	// delete the invite from GroupInvite
	_, err := db.SQL.Exec(`DELETE FROM "GroupInvite" WHERE userID = ? AND groupID = ?;`, userID, groupID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// LeaveGroup -
func LeaveGroup(userID string, groupID string) error {

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if user is group owner
	// check if user exists in group
	var userIsOwner bool
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE id = ? AND ownerID = ?),
						EXISTS(SELECT 1 FROM "UserGroup" WHERE groupID = ? AND userID = ?);`, groupID, userID, groupID, userID).Scan(&userIsOwner, &userExistsInGroup)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// if user is owner
	if userIsOwner {
		return errors.New("user is owner")
	}

	// if user is not in group
	if !userExistsInGroup {
		return errors.New("user not in group")
	}

	err = deleteUserGroup(tx, userID, groupID)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// RemoveUserFromGroup -
func RemoveUserFromGroup(ownerID string, userID string, groupID string) error {

	if ownerID == userID {
		return errors.New("invalid user")
	}

	// check if group owner
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if user is group owner
	// check if user exists in group
	var validOwner bool
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE id = ? AND ownerID = ?),
						EXISTS(SELECT 1 FROM "UserGroup" WHERE groupID = ? AND userID = ?);`,
		groupID, ownerID, groupID, userID).Scan(&validOwner, &userExistsInGroup)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// if user is owner
	if !validOwner {
		return errors.New("invalid owner")
	}

	// if user is not in group
	if !userExistsInGroup {
		return errors.New("user not in group")
	}

	// delete from UserGroup where userID and groupID
	err = deleteUserGroup(tx, userID, groupID)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// DeleteGroup -
func DeleteGroup(ownerID string, groupID string) error {

	// check if group owner
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if user is group owner
	var validOwner bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE id = ? AND ownerID = ?);`, groupID, ownerID).Scan(&validOwner)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// if user is owner
	if !validOwner {
		return errors.New("invalid owner")
	}

	_, err = tx.Exec(`DELETE FROM "Message" WHERE groupID = ?;
						DELETE FROM "GeoLocation" WHERE groupID = ?;
						DELETE FROM "GroupInvite" WHERE groupID = ?;
						DELETE FROM "UserGroup" WHERE groupID = ?;
						DELETE FROM "Group" WHERE id = ?;`,
		groupID, groupID, groupID, groupID, groupID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

func insertUserGroup(tx *sql.Tx, userID string, groupID string) error {

	// insert id's into UserGroup table
	_, err := tx.Exec(`INSERT INTO "UserGroup" (userID, groupID) SELECT ?, ?
					WHERE NOT EXISTS(SELECT 1 FROM "UserGroup" WHERE userID = ? AND groupID = ?);`, userID, groupID, userID, groupID)

	if err != nil {
		log.Println(err)
		return err
	}

	// update userCount in Group table
	_, err = tx.Exec(`UPDATE "Group" SET userCount = userCount + 1 WHERE id = ?;`, groupID)

	if err != nil {
		log.Println(err)
	}

	return err
}

func deleteUserGroup(tx *sql.Tx, userID string, groupID string) error {
	// delete from UserGroup where userID and groupID
	_, err := tx.Exec(`DELETE FROM "UserGroup" WHERE userID = ? AND groupID = ?;
						DELETE FROM "GeoLocation" WHERE userID = ? AND groupID = ?;`, userID, groupID, userID, groupID)

	if err != nil {
		return err
	}

	// update userCount in Group table
	_, err = tx.Exec(`UPDATE "Group" SET userCount = userCount - 1 WHERE id = ?;`, groupID)

	return err
}

// UpdateGroupInfo - update group to public/private and update password
func UpdateGroupInfo(ownerID string, groupID string, password string, public bool) error {

	// start db transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	// check if user is group owner
	// check if user exists in group
	var validOwner bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "Group" WHERE id = ? AND ownerID = ?);`,
		groupID, ownerID).Scan(&validOwner)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// if user is owner
	if !validOwner {
		return errors.New("Invalid group owner.")
	}

	// check if the password is updated
	if password != "" {

		// hash the password before storing in the database
		passwordHash, err := utils.GenerateHash(password)

		if err != nil {
			log.Println(err)
			return errors.New("Error hashing password")
		}

		_, err = tx.Exec(`UPDATE "Group" SET "public" = ?, "password" = ? WHERE "id" = ?;`, public, passwordHash, groupID)
	} else {

		// set password to null if password is blank
		_, err = tx.Exec(`UPDATE "Group" SET "public" = ?, "password" = NULL WHERE "id" = ?;`, public, groupID)
	}

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}
