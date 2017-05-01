package groupOperations

import (
	"errors"
	"log"

	fcm "github.com/NaySoftware/go-fcm"
	"github.com/mgerb/wbu-server/config"
	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/model"
)

//StoreUserGroupMessages - store a users messages in a group
func StoreUserGroupMessages(groupID string, userID string, message string) error {

	//do validation before running redis script
	//check message length - must be less than 150 characters
	if len(message) > 150 || len(message) == 0 {
		return errors.New("invalid message")
	}

	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

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

	// must commit before starting new db transaction
	tx.Commit()

	// send out notifications via FCM - new go routine
	go fcmNotifications(groupID, userID, message)

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
	rows, err := tx.Query(`SELECT m.id, m.content, m.userID, m.groupID, u.firstName, u.lastName, strftime('%s', m.timestamp) FROM "Message" AS m INNER JOIN
	   						"User" AS u ON u.id = m.userID
	   						WHERE m.groupID = ? AND datetime(m.timestamp, 'localtime') > datetime(?, 'unixepoch', 'localtime');`, groupID, unixTime) //timestamp.Format("2006-01-02 15:04:05"))

	if err != nil {
		log.Println(err)
		return []*model.Message{}, errors.New("database error")
	}

	defer rows.Close()

	messageList := []*model.Message{}

	for rows.Next() {
		newMessage := &model.Message{}
		err := rows.Scan(&newMessage.ID, &newMessage.Content, &newMessage.UserID, &newMessage.GroupID, &newMessage.FirstName, &newMessage.LastName, &newMessage.Timestamp)

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

// TODO rethink this - need to get group name along with users name
// fcmNotifications - get all fcm tokens and send notifications to FCM
func fcmNotifications(groupID string, userID string, message string) {

	userList, err := db.RClient.SMembers(model.UserGroupKey + groupID).Result()

	if err != nil {
		log.Println(err)
		return
	}

	tokenList, err := db.RClient.HMGet(model.FCMTokenKey, userList...).Result()

	if err != nil {
		log.Println(err)
		return
	}

	var tokenStringList []string
	for _, token := range tokenList {
		tokenStringList = append(tokenStringList, token.(string))
	}

	data := map[string]interface{}{
		"msg": "Test 123",
	}

	client := fcm.NewFcmClient(config.Config.FCMServerKey)

	notifPayload := &fcm.NotificationPayload{
		Title: "Group Name",
		Body:  message,
		Sound: "Enabled",
	}

	client.SetNotificationPayload(notifPayload)

	client.NewFcmRegIdsMsg(tokenStringList, data)

	status, err := client.Send()

	if err != nil {
		log.Println(err)
		return
	}

	status.PrintResults()
}
