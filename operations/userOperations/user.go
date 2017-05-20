package userOperations

import (
	"errors"
	"log"
	"regexp"
	"strconv"

	"database/sql"

	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/model"
	"github.com/mgerb/wbu-server/operations/fb"
	"github.com/mgerb/wbu-server/utils"
	"github.com/mgerb/wbu-server/utils/regex"
	"github.com/mgerb/wbu-server/utils/tokens"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser -
func CreateUser(email string, password string, firstName string, lastName string) error {

	//validate password
	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("Password must be at least 5 characters.")
	}

	passwordHash, err := utils.GenerateHash(password)

	if err != nil {
		log.Println(err)
		return errors.New("Internal error.")
	}

	//validate email
	if !regexp.MustCompile(regex.EMAIL).MatchString(email) {
		return errors.New("Invalid email.")
	}

	nameLength := len(firstName + lastName)

	// validate first/last name
	if nameLength < 1 || nameLength > 40 {
		return errors.New("Invalid name.")
	}

	// start sql transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("Internal error.")
	}

	// commit the transaction when the function returns
	defer tx.Commit()

	//check if the email already exists
	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "User" WHERE "email" = ?);`, email).Scan(&userExists)

	// return err or if email already exists
	if err != nil {
		log.Println(err)
		return errors.New("Internal error.")
	} else if userExists {
		return errors.New("Email is taken.")
	}

	err = insertUserInformation(tx, email, passwordHash, firstName, lastName, "")

	if err != nil {
		log.Println(err)
		return errors.New("Internal error.")
	}

	return nil
}

func insertUserInformation(tx *sql.Tx, email string, passwordHash string, firstName string, lastName string, facebookID string) error {

	// if not facebook user
	if facebookID != "" {
		_, err := tx.Exec(`INSERT INTO "User" (email, firstName, lastName, facebookID) VALUES (?, ?, ?, ?);`, email, firstName, lastName, facebookID)

		if err != nil {
			return err
		}

	} else {
		// insert into User email, passwordHash, firstName, and lastName
		_, err := tx.Exec(`INSERT INTO "User" (email, password, firstName, lastName) VALUES (?, ?, ?, ?);`, email, passwordHash, firstName, lastName)

		if err != nil {
			return err
		}
	}

	// populate UserSettings table
	_, err := tx.Exec(`INSERT INTO "UserSettings" (userID) VALUES((SELECT id FROM "User" WHERE email = ?));`, email)

	if err != nil {
		return err
	}

	return nil
}

//Login - check if valid credentials - create JWT and return User object
func Login(email string, password string) (*model.User, error) {

	// new user struct
	newUser := &model.User{}

	// get user information
	err := db.SQL.QueryRow(`SELECT id, email, firstName, lastName, password FROM "User" WHERE "email" = ?;`, email).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.Password)

	if err != nil {
		log.Println(err)
		return newUser, errors.New("Invalid login.")
	}

	if bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password)) != nil {
		return &model.User{}, errors.New("Invalid login.")
	}

	token, lastRefreshTime, errToken := tokens.GetJWT(newUser.Email, strconv.FormatInt(newUser.ID, 10), newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &model.User{}, errors.New("Internal error.")
	}

	newUser.Jwt = token
	newUser.LastRefreshTime = lastRefreshTime
	newUser.FacebookUser = false

	return newUser, nil
}

// LoginFacebook -
func LoginFacebook(accessToken string) (*model.User, error) {
	// check if valid facebook user
	fbResponse, err := fb.Me(accessToken)
	if err != nil {
		return &model.User{}, errors.New("Internal error.")
	}

	// get the facebook user id
	facebookID := fbResponse["id"].(string)

	// new user struct
	newUser := &model.User{}

	// get user information
	err = db.SQL.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID)

	// store user info if not exists
	if err == sql.ErrNoRows {
		// start sql transaction
		tx, err := db.SQL.Begin()
		if err != nil {
			log.Println(err)
			return &model.User{}, errors.New("Internal error.")
		}

		// commit the transaction when the function returns
		defer tx.Commit()

		// get facebook user's information from token
		email := fbResponse["email"].(string)
		firstName := fbResponse["first_name"].(string)
		lastName := fbResponse["last_name"].(string)
		facebookID := fbResponse["id"].(string)

		err = insertUserInformation(tx, email, "", firstName, lastName, facebookID)

		if err != nil {
			log.Println(err)
			return &model.User{}, errors.New("Internal error.")
		}

		// get user information
		err = tx.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID.String)

		if err != nil {
			log.Println(err)
			return &model.User{}, errors.New("Internal error.")
		}

	} else if err != nil {
		log.Println(err)
		return &model.User{}, errors.New("Internal error.")
	}

	// generate new json web token
	token, lastRefreshTime, errToken := tokens.GetJWT(newUser.Email, strconv.FormatInt(newUser.ID, 10), newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &model.User{}, errors.New("Internal error.")
	}

	newUser.Jwt = token
	newUser.LastRefreshTime = lastRefreshTime
	newUser.FacebookUser = true

	return newUser, nil
}

// SearchUserByName - return list of users that match name
func SearchUserByName(searchTerm string, userID string) ([]*model.User, error) {

	userList := []*model.User{}

	// get user information
	rows, err := db.SQL.Query(`SELECT id, email, firstName, lastName FROM "User" WHERE "firstName" || ' ' || "lastName" LIKE ? OR "email" LIKE ? AND "id" != ? LIMIT 20;`, "%"+searchTerm+"%", "%"+searchTerm+"%", userID)

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("Internal error.")
	}

	defer rows.Close()

	// map query to user object list
	for rows.Next() {
		newUser := &model.User{}
		err := rows.Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName)

		if err != nil {
			log.Println(err)
			return []*model.User{}, errors.New("Internal error.")
		}
		userList = append(userList, newUser)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("Internal error.")
	}

	return userList, nil
}

// UpdateFCMToken - Update fcm token in user table
func UpdateFCMToken(userID string, token string) error {

	if token == "" {
		return errors.New("invalid token")
	}

	_, err := db.SQL.Exec(`UPDATE "UserSettings" SET fcmToken = ? WHERE userID = ?;`, token, userID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// ToggleNotifications - Update fcm token in user table
func ToggleNotifications(userID string, toggle string) error {

	if toggle != "1" && toggle != "0" {
		return errors.New("invalid value")
	}

	_, err := db.SQL.Exec(`UPDATE "UserSettings" SET notifications = ? WHERE userID = ?;`, toggle, userID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// GetUserSettings - get user settings
func GetUserSettings(userID string) (*model.UserSettings, error) {

	settings := &model.UserSettings{}

	err := db.SQL.QueryRow(`SELECT userID, notifications FROM "UserSettings" WHERE userID = ?;`, userID).Scan(&settings.UserID, &settings.Notifications)

	if err != nil {
		log.Println(err)
		return &model.UserSettings{}, err
	}

	return settings, nil
}

// DeleteUser -
func DeleteUser(userID string) error {

	// check if user exists

	// delete from UserGroup table

	// delete from User table

	return nil
}

// RemoveFCMToken -
func RemoveFCMToken(userID string) error {

	_, err := db.SQL.Exec(`UPDATE "UserSettings" SET fcmToken = NULL WHERE userID = ?;`, userID)

	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	return nil
}

// StoreUserFeedback -
func StoreUserFeedback(userID string, feedback string) error {

	if feedback == "" {
		return errors.New("Invalid feedback.")
	}

	_, err := db.SQL.Exec(`INSERT INTO "UserFeedback"(userID, feedback) VALUES(?,?);`, userID, feedback)

	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	return nil
}
