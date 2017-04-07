package groupOperations

import (
	"errors"
	"log"

	"../../db"
	"../../model"
)

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

	// send out notifications via FCM
	fcmNotifications(groupID, userID)

	return nil
}

//GetUserGroupMessages - return user messages for a group
func GetUserGroupMessages(groupID string, userID string, unixTime string) ([]*model.Message, error) {

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

	// need to convert time input to local time because sqlite compares time strings and not unix time
	rows, err := tx.Query(`SELECT m.id, m.content, u.firstName, u.lastName, strftime('%s', m.timestamp) FROM "Message" AS m INNER JOIN
	   						"User" AS u ON u.id = m.userID
	   						WHERE m.groupID = ? AND datetime(m.timestamp, 'localtime') >= datetime(?, 'unixepoch', 'localtime');`, groupID, unixTime) //timestamp.Format("2006-01-02 15:04:05"))

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	}

	defer rows.Close()

	messageList := []*model.Message{}

	for rows.Next() {
		newMessage := &model.Message{}
		err := rows.Scan(&newMessage.ID, &newMessage.Content, &newMessage.FirstName, &newMessage.LastName, &newMessage.Timestamp)

		if err != nil {
			log.Println(err)
			return []*model.Message{}, errors.New("database error")
		}

		messageList = append(messageList, newMessage)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	}

	return messageList, nil
}

// fcmNotifications - get all fcm tokens and send notifications to FCM exclude userID from notifications
func fcmNotifications(groupID string, userID string) {

	rows, err := db.SQL.Query(`SELECT u.fcmToken from "UserGroup" AS ug INNER JOIN
								"User" AS u ON ug.userID = u.id WHERE ug.groupID = ? AND ug.userID != ? AND u.fcmToken != NULL;`, groupID, userID)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	tokenList := []string{}

	for rows.Next() {
		var token string
		err := rows.Scan(&token)

		if err != nil {
			log.Println(err)
		}

		tokenList = append(tokenList, token)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
	}

	// TODO - add FCM functionality
	log.Println("FCM Token List")
	log.Println(tokenList)
}
