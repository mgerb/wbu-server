package userOperations

import (
	"errors"
	"log"
	"regexp"

	"../lua"
	redis "gopkg.in/redis.v5"

	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils/regex"
	"../../utils/tokens"
	"../fb"
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

	//validate full name
	/* TODO Validate names
	if !regexp.MustCompile(regex.FULL_NAME).MatchString(firstName + " " + lastName) {
		return errors.New("Invalid name.")
	}
	*/

	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("Database error.")
	}

	defer tx.Commit()

	//check if the email already exists
	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "User" WHERE email = ?);`, email).Scan(&userExists)

	if err != nil || userExists {
		return errors.New("User already exists")
	}
	log.Println(lastName)
	_, err = tx.Exec(`INSERT INTO "User" (email, password, firstName, lastName) VALUES (?, ?, ?, ?);`, email, passwordHash, firstName, lastName, "test")

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
	err := db.SQL.QueryRow(`SELECT id, email, firstName, lastName, password FROM "User" WHERE EMAIL = ?;`, email).Scan(&newUser.UserID, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.Password)

	if err != nil {
		log.Println(err)
		return newUser, errors.New("User does not exist.")
	}

	if bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password)) != nil {
		return &userModel.User{}, errors.New("Invalid password.")
	}

	token, lastJwtRefresh, errToken := tokens.GetJWT(newUser.Email, newUser.UserID, newUser.FirstName, newUser.LastName)

	if errToken != nil {
		log.Println(err)
		return &userModel.User{}, errors.New("Token error.")
	}

	newUser.Jwt = token
	newUser.LastJwtRefresh = lastJwtRefresh

	return newUser, nil
}

func LoginFacebook(accessToken string) (map[string]interface{}, error) {
	response, err_fb := fb.Me(accessToken)

	if err_fb != nil {
		return map[string]interface{}{}, errors.New("Invalid FB token")
	}

	email := response["email"].(string)
	fullName := response["name"].(string)
	facebookID := response["id"].(string)

	//prepend facebookID with "fb:" to not get mixed with regular user id's
	userID := userModel.GetUserID_FB(facebookID)

	//check if hash key exists
	userExists := db.Client.Exists(userModel.USER_HASH(userID)).Val()

	//create user if not exists
	if !userExists {
		pipe := db.Client.Pipeline()
		defer pipe.Close()

		//map users email to new id
		pipe.HSet(userModel.USER_ID(), email, userID)

		//set user object in redis
		pipe.HMSet(userModel.USER_HASH(userID), userModel.USER_HASH_MAP(email, "", "0", fullName))

		_, err_pipe := pipe.Exec()

		if err_pipe != nil {
			return map[string]interface{}{}, errors.New("pipe error")
		}
	}

	token, lastRefreshTime, err_token := tokens.GetJWT(email, userID, fullName, "TODO REMOVE THIS")

	return map[string]interface{}{
		"email":           email,
		"userID":          userID,
		"fullName":        fullName,
		"jwt":             token,
		"lastRefreshTime": lastRefreshTime,
	}, err_token
}

func DeleteUser(userID string) error {

	script := redis.NewScript(lua.Use("DeleteUser.lua"))

	return script.Run(db.Client, []string{
		userModel.USER_HASH(userID),
		userModel.USER_ID(),
		userModel.USER_GROUPS(userID),
		userModel.USER_GROUP_INVITES(userID),
	},
		userID,
		userModel.USER_GROUP_MESSAGES_KEY(),
		groupModel.GROUP_MEMBERS_KEY(),
		groupModel.GROUP_LOCATIONS_KEY(),
	).Err()
}

//GetUserGroups - get all the groups the user exists in
func GetGroups(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUPS(userID)).Result()
}

func GetInvites(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUP_INVITES(userID)).Result()
}

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
