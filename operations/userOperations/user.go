package userOperations

import (
	"errors"
	"log"
	"regexp"

	"../../db"
	"../../model/userModel"
	"../../utils/regex"
	"../../utils/tokens"
	"../fb"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

//CreateUser - store userName/password in hash
func CreateUser(email string, password string, firstName string, lastName string) error {

	//validate password
	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("Invalid password.")
	}

	passwordHash := generateHash(password)

	//validate email
	if !regexp.MustCompile(regex.EMAIL).MatchString(email) {
		return errors.New("Invalid email.")
	}

	if len(firstName+lastName) > 40 {
		return errors.New("Name is too long")
	}

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	defer tx.Commit()

	//check if the email already exists
	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "User" WHERE "email" = ?);`, email).Scan(&userExists)

	if err != nil || userExists {
		return errors.New("User already exists")
	}

	_, err = tx.Exec(`INSERT INTO "User" (email, password, firstName, lastName) VALUES (?, ?, ?, ?);`, email, passwordHash, firstName, lastName)

	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	return nil
}

//Login - check if password and userName are correct
func Login(email string, password string) (*userModel.User, error) {

	// new user struct
	newUser := &userModel.User{}

	// get user information
	err := db.SQL.QueryRow(`SELECT id, email, firstName, lastName, password FROM "User" WHERE "email" = ?;`, email).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.Password)

	if err != nil {
		log.Println(err)
		return newUser, errors.New("User does not exist.")
	}

	if bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password)) != nil {
		return &userModel.User{}, errors.New("Invalid password.")
	}

	token, lastJwtRefresh, errToken := tokens.GetJWT(newUser.Email, newUser.ID, newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &userModel.User{}, errors.New("Token error.")
	}

	newUser.Jwt = token
	newUser.LastJwtRefresh = lastJwtRefresh

	return newUser, nil
}

func LoginFacebook(accessToken string) (*userModel.User, error) {
	// check if valid facebook user
	fbResponse, err := fb.Me(accessToken)
	if err != nil {
		return &userModel.User{}, errors.New("Invalid FB token")
	}

	// get the facebook user id
	facebookID := fbResponse["id"].(string)

	// new user struct
	newUser := &userModel.User{}

	// get user information
	err = db.SQL.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID)

	// store user info if not exists
	if err == sql.ErrNoRows {
		err := createFacebookUser(fbResponse)

		if err != nil {
			log.Println(err)
			return &userModel.User{}, errors.New("Database error.")
		}

		// get user information
		err = db.SQL.QueryRow(`SELECT id, email, firstName, lastName, facebookID FROM "User" WHERE "facebookID" = ?;`, facebookID).Scan(&newUser.ID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.FacebookID)

		if err != nil {
			log.Println(err)
			return &userModel.User{}, errors.New("Database error.")
		}

	} else if err != nil {
		log.Println(err)
		return &userModel.User{}, errors.New("Database error.")
	}

	token, lastJwtRefresh, errToken := tokens.GetJWT(newUser.Email, newUser.ID, newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &userModel.User{}, errors.New("Token error.")
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

func DeleteUser(userID string) error {
	return nil
}

//GetUserGroups - get all the groups the user exists in
func GetGroups(userID string) (map[string]string, error) {
	return map[string]string{"test": "123"}, nil
}

func GetInvites(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUP_INVITES(userID)).Result()
}

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
