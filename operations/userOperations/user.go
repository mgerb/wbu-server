package userOperations

import (
	"errors"
	"strconv"
	"time"

	"../../config"
	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"../fb"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	//expiration time for JWT
	expirationTime int64 = 30 * 24 * 60 * 60
)

/*
//CreateUser - store userName/password in hash
func CreateUser(userName string, password string) error {

	//DO VALIDATION
	if !regexp.MustCompile(regex.USERNAME).MatchString(userName) {
		return errors.New("Invalid userName.")
	}

	if !regexp.MustCompile(regex.PASSWORD).MatchString(password) {
		return errors.New("Invalid password.")
	}

	//check if the userName already exists in redis
	_, err := GetUserID(userName)
	if err == nil {
		return errors.New("userName already exists.")
	}

	temp, _ := db.Client.Incr(userModel.USER_KEY_STORE()).Result()
	newID := strconv.FormatInt(temp, 10)

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.Set(userModel.USER_ID(userName), newID, 0)

	//set user object in redis
	pipe.HMSet(userModel.USER_HASH(newID), userModel.USER_HASH_MAP(userName, generateHash(password), "0"))

	_, err = pipe.Exec()

	return err
}

//Login - check if password and userName are correct
func Login(userName string, password string) (string, error) {
	id, err := GetUserID(userName)

	if err == nil {
		result, _ := db.Client.HGet(userModel.USER_HASH(id), "password").Result()
		if bcrypt.CompareHashAndPassword([]byte(result), []byte(password)) != nil {
			return "", errors.New("Invalid password.")
		}
	} else {
		return "", errors.New("User does not exist.")
	}

	//if user has valid login - generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": userName,
		"userID":   id,
		"exp":      time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))

	if tokenError != nil {
		return "", errors.New("Token error.")
	}

	return tokenString, nil
}
*/

func LoginFacebook(accessToken string) (string, error) {
	response, err_fb := fb.Me(accessToken)

	if err_fb != nil {
		return "", errors.New("Invalid FB token")
	}

	email := response["email"].(string)
	usersName := response["name"].(string)
	facebookID := response["id"].(string)

	userID, exists_err := db.Client.HGet(userModel.USER_ID(), email).Result()

	//create user if not exists
	if exists_err != nil {
		temp, err_incr := db.Client.Incr(userModel.USER_KEY_STORE()).Result()

		if err_incr != nil {
			return "", errors.New("Incr error")
		}

		userID = strconv.FormatInt(temp, 10)

		pipe := db.Client.Pipeline()
		defer pipe.Close()

		pipe.HSet(userModel.USER_ID(), email, userID)

		//set user object in redis
		pipe.HMSet(userModel.USER_HASH(userID), userModel.USER_HASH_MAP("", "", "0", usersName, email, facebookID))

		_, err_pipe := pipe.Exec()

		if err_pipe != nil {
			return "", errors.New("pipe error")
		}
	}

	//log user in
	//if user has valid login - generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"userID":    userID,
		"usersName": usersName,
		"exp":       time.Now().Unix() + expirationTime,
	})

	tokenString, tokenError := token.SignedString([]byte(config.Config.TokenSecret))

	if tokenError != nil {
		return "", errors.New("Token error.")
	}

	return tokenString, nil
}

//GetUserID = return user id as string
func GetUserID(email string) (string, error) {
	return db.Client.HGet(userModel.USER_ID(), email).Result()
}

//GetUserGroups - get all the groups the user exists in
func GetGroups(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUPS(userID)).Result()
}

func GetInvites(userID string) (map[string]string, error) {
	return db.Client.HGetAll(userModel.USER_GROUP_INVITES(userID)).Result()
}

//TODO----------------------------------------------------------
func JoinGroup(userID string, name string, groupID string) error {

	userHasInvite := db.Client.HExists(userModel.USER_GROUP_INVITES(userID), groupID).Val()

	if !userHasInvite {
		return errors.New("User does not have an invite.")
	}

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	pipe.HSet(groupModel.GROUP_MEMBERS(groupID), userID, name)
	pipe.HDel(userModel.USER_GROUP_INVITES(userID), groupID)

	_, err := pipe.Exec()

	if err != nil {
		return err
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
