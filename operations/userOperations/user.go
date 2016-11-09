package userOperations

import (
	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../../utils/regex"
	"../../utils/tokens"
	"../fb"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
)

//CreateUser - store userName/password in hash
func CreateUser(email string, password string, fullName string) error {

	//TODO: validation for email
	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("Invalid password.")
	}

	//check if the email already exists in redis
	emailExists := db.Client.HExists(userModel.USER_ID(), email).Val()

	if emailExists {
		return errors.New("That email is already taken.")
	}

	//get a new user id from the id pool
	temp, _ := db.Client.Incr(userModel.USER_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	//prepend user id with prefix
	newID = userModel.GetUserID(newID)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//map email to user id
	pipe.HSet(userModel.USER_ID(), email, newID)

	//set user object in redis
	pipe.HMSet(userModel.USER_HASH(newID), userModel.USER_HASH_MAP(email, generateHash(password), "0", fullName))

	_, err := pipe.Exec()

	return err
}

//Login - check if password and userName are correct
func Login(email string, password string) (string, error) {
	userID, err := db.Client.HGet(userModel.USER_ID(), email).Result()

	if err != nil {
		return "", errors.New("User does not exist.")
	}

	result, err_password := db.Client.HGetAll(userModel.USER_HASH(userID)).Result()

	if err_password != nil {
		return "", errors.New("Error retrieving account information.")
	}

	savedPassword := result["password"]
	fullName := result["fullName"]

	if len(savedPassword) < 5 {
		return "", errors.New("Invalid account")
	}

	if bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password)) != nil {
		return "", errors.New("Invalid password.")
	}

	return tokens.GetJWT(email, userID, fullName)
}

func LoginFacebook(accessToken string) (string, error) {
	response, err_fb := fb.Me(accessToken)

	if err_fb != nil {
		return "", errors.New("Invalid FB token")
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
			return "", errors.New("pipe error")
		}
	}

	return tokens.GetJWT(email, userID, fullName)
}

//GetUserGroups - get all the groups the user exists in
func GetGroups(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUPS(userID)).Result()
}

func GetInvites(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUP_INVITES(userID)).Result()
}

//TODO----------------------------------------------------------
func JoinGroup(userID string, groupID string) error {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	userInfo := pipe.HGetAll(userModel.USER_HASH(userID))
	groupInfo := pipe.HGetAll(groupModel.GROUP_HASH(groupID))

	userHasInvite := pipe.HExists(userModel.USER_GROUP_INVITES(userID), groupID)

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return err_pipe1
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

func LeaveGroup(userID string, groupid string) error {
	return errors.New("TODO")
}

//TODO----------------------------------------------------------

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}
