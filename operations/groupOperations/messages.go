package groupOperations

import (
	"errors"
	"log"

	"../../db"
	"../../model"
)

// TODO - test
//StoreUserGroupMessages - store a users messages in a group
func StoreUserGroupMessages(groupID string, userID string, message string) error {

	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	//do validation before running redis script
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("invalid message")
	}

	// check if user exists in group before sending messages
	// check ID's in UserGroup table
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE "groupID" = ? AND "userID" = ?);`, groupID, userID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	} else if !userExistsInGroup {
		return errors.New("user not in group")
	}

	_, err = tx.Exec(`INSERT INTO "Message" (groupID, userID, content) VALUES (?, ?, ?);`, groupID, userID, message)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// TODO - TEST
//GetUserGroupMessages - return user messages for a group
func GetUserGroupMessages(groupID string, userID string, timestamp string) ([]*model.Message, error) {

	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	}

	defer tx.Commit()

	// check if user exists in group before sending messages
	// check ID's in UserGroup table
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE "groupID" = ? AND "userID" = ?);`, groupID, userID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	} else if !userExistsInGroup {
		return []*model.Message{}, errors.New("user not in group")
	}

	rows, err := tx.Query(`SELECT m.id, m.content, u.firstName, u.lastName FROM "Message" AS m INNER JOIN
                        "UserGroup" AS ug ON m.groupID = ug.groupID INNER JOIN
						"User" AS u ON u.id = ug.userID
						WHERE groupID = ? AND timestamp > ?;`, groupID, timestamp)

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	}

	defer rows.Close()

	messageList := []*model.Message{}

	for rows.Next() {
		newMessage := &model.Message{}
		err := rows.Scan(&newMessage.ID, &newMessage.Content, &newMessage.FirstName, &newMessage.LastName)

		if err != nil {
			log.Println(err)
			return []*model.Message{}, errors.New("database error")
		}

		messageList = append(messageList, newMessage)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("row error")
	}

	return []*model.Message{}, nil
}

// fcmNotifications - get all fcm tokens and
func fcmNotifications(groupID string) {

}
