package userOperations

import (
	"errors"
	"log"
	"regexp"

	"database/sql"

	"../../db"
	"../../model"
	"../../utils"
	"../../utils/regex"
	"../../utils/tokens"
	"../fb"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser -
func CreateUser(email string, password string, firstName string, lastName string) error {

	//validate password
	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("invalid password")
	}

	passwordHash, err := utils.GenerateHash(password)

	if err != nil {
		log.Println(err)
		return errors.New("password hash error")
	}

	//validate email
	if !regexp.MustCompile(regex.EMAIL).MatchString(email) {
		return errors.New("invalid email")
	}

	nameLength := len(firstName + lastName)

	// validate first/last name
	if nameLength < 2 || nameLength > 40 {
		return errors.New("invalid name")
	}

	// start sql transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	// commit the transaction when the function returns
	defer tx.Commit()

	//check if the email already exists
	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "User" WHERE "email" = ?);`, email).Scan(&userExists)

	// return err or if email already exists
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	} else if userExists {
		return errors.New("email taken")
	}

	// insert into User email, passwordHash, firstName, and lastName
	_, err = tx.Exec(`INSERT INTO "User" (email, password, firstName, lastName) VALUES (?, ?, ?, ?);`, email, passwordHash, firstName, lastName)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
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
		return newUser, errors.New("invalid user")
	}

	if bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password)) != nil {
		return &model.User{}, errors.New("invalid password")
	}

	token, lastJwtRefresh, errToken := tokens.GetJWT(newUser.Email, newUser.ID, newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &model.User{}, errors.New("jwt error")
	}

	newUser.Jwt = token
	newUser.LastJwtRefresh = lastJwtRefresh

	return newUser, nil
}

// LoginFacebook -
func LoginFacebook(accessToken string) (*model.User, error) {
	// check if valid facebook user
	fbResponse, err := fb.Me(accessToken)
	if err != nil {
		return &model.User{}, errors.New("facebook error")
	}

	// get the facebook user id
	facebookID := fbResponse["id"].(string)

	// new user struct
	newUser := &model.User{}

	// get user information
	err = db.SQL.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID)

	// store user info if not exists
	if err == sql.ErrNoRows {
		err := createFacebookUser(fbResponse)

		if err != nil {
			log.Println(err)
			return &model.User{}, errors.New("database error")
		}

		// get user information
		err = db.SQL.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID.String)

		if err != nil {
			log.Println(err)
			return &model.User{}, errors.New("database error")
		}

	} else if err != nil {
		log.Println(err)
		return &model.User{}, errors.New("database error")
	}

	token, lastJwtRefresh, errToken := tokens.GetJWT(newUser.Email, newUser.ID, newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &model.User{}, errors.New("jwt error")
	}

	newUser.Jwt = token
	newUser.LastJwtRefresh = lastJwtRefresh

	return newUser, nil
}

func createFacebookUser(fbResponse map[string]interface{}) error {
	// get facebook user's information from token
	email := fbResponse["email"].(string)
	firstName := fbResponse["first_name"].(string)
	lastName := fbResponse["last_name"].(string)
	facebookID := fbResponse["id"].(string)

	_, err := db.SQL.Exec(`INSERT INTO "User" (email, firstName, lastName, facebookID) VALUES (?, ?, ?, ?);`, email, firstName, lastName, facebookID)

	if err != nil {
		return err
	}

	return nil
}

// SearchUserByName - return list of users that match name
func SearchUserByName(name string) ([]*model.User, error) {

	userList := []*model.User{}

	// get user information
	rows, err := db.SQL.Query(`SELECT id, email, firstName, lastName FROM "User" WHERE "firstName" || ' ' || "lastName" LIKE ?;`, "%"+name+"%")

	if err != nil {
		log.Println(err)
		return []*model.User{}, errors.New("database error")
	}

	defer rows.Close()

	// map query to user object list
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

// DeleteUser -
func DeleteUser(userID string) error {

	// check if user exists

	// delete from UserGroup table

	// delete from User table

	return nil
}

// UpdateFCMToken - Update fcm token in user table
func UpdateFCMToken(userID string, token string) error {

	if token == "" {
		return errors.New("invalid token")
	}

	_, err := db.SQL.Exec(`UPDATE "User" SET fcmToken = ? WHERE id = ?;`, token, userID)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}
