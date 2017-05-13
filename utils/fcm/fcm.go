package fcm

import (
	"log"

	fcm "github.com/NaySoftware/go-fcm"
	"github.com/mgerb/wbu-server/config"
	"github.com/mgerb/wbu-server/db"
)

type data struct {
	NotifType string `json:"type"`
	GroupID   string `json:"groupID,omitempty"`
}

// SendToGroup - sends message to a group that the user is in
func SendToGroup(groupID string, senderUserID string, message string, notifType string) {
	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	defer tx.Commit()

	// get users name
	var firstName, lastName string
	err = tx.QueryRow(`SELECT firstName, lastName from "User" WHERE id = ?;`, senderUserID).Scan(&firstName, &lastName)

	if err != nil {
		log.Println(err)
		return
	}

	rows, err := tx.Query(`SELECT us.fcmToken FROM "Group" AS g
								INNER JOIN "UserGroup" AS ug ON g.id = ug.groupID
								INNER JOIN "UserSettings" AS us ON ug.userID = us.userID
								WHERE g.id = ? AND ug.userID != ? AND us.notifications = 1;`, groupID, senderUserID)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	tokenList := []string{}

	for rows.Next() {
		var token string
		err := rows.Scan(&token)
		if err != nil {
			log.Println(err)
			continue
		}

		tokenList = append(tokenList, token)
	}

	// set message title to users name
	title := firstName + " " + lastName

	// create data object
	data := &data{
		NotifType: notifType,
		GroupID:   groupID,
	}

	err = sendNotif(title, message, data, tokenList)

	if err != nil {
		log.Println(err)
	}
}

// UserInviteNotif -
func UserInviteNotif(userID string) {
	var fcmToken string
	err := db.SQL.QueryRow(`SELECT fcmToken FROM "UserSettings" WHERE userID = ? AND notifications == 1;`, userID).Scan(&fcmToken)

	if err != nil {
		log.Println(err)
		return
	}

	title := "Group Invite"
	body := "You have a new group invite!"
	data := &data{
		NotifType: "groupInvite",
	}

	sendNotif(title, body, data, []string{fcmToken})
}

func sendNotif(title string, body string, data *data, tokenList []string) error {

	client := fcm.NewFcmClient(config.Config.FCMServerKey)

	notifPayload := &fcm.NotificationPayload{
		Title: title,
		Body:  body,
		Sound: "Enabled",
	}

	client.SetNotificationPayload(notifPayload)

	client.NewFcmRegIdsMsg(tokenList, data)

	_, err := client.Send()

	return err
}
