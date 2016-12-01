package groupOperations

import (
	"../../db"
	"../../model/groupModel"
	"../../model/userModel"
	"errors"
	redis "gopkg.in/redis.v5"
	"strconv"
	"time"
)

func StoreGeoLocation(userID string, groupID string, latitude string, longitude string) error {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//DO VALIDATION
	//check if user exists in group before storing message
	userIsMember := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), userID)
	fullName := pipe.HGet(userModel.USER_HASH(userID), "fullName")

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return errors.New("Pipe error.")
	}

	if !userIsMember.Val() {
		return errors.New("You are not a member of this group.")
	}

	long, err1 := strconv.ParseFloat(longitude, 64)
	lat, err2 := strconv.ParseFloat(latitude, 64)

	if err1 != nil || err2 != nil {
		return errors.New("Invalid longitude/latitude.")
	}

	//validate coordinates
	if long < -180 ||
		long > 180 ||
		lat < -85.05112878 ||
		lat > 85.05112878 {
		return errors.New("Invalid coordinates.")
	}

	geoLocation := &redis.GeoLocation{}

	geoLocation.Name = userID + "/" + fullName.Val() + "/" + strconv.FormatInt(time.Now().Unix(), 10)
	geoLocation.Longitude = long
	geoLocation.Latitude = lat

	_, returnError := db.Client.GeoAdd(groupModel.GROUP_GEO(groupID), geoLocation).Result()

	return returnError
}

func GetGeoLocations(groupID string, userID string) ([]string, error) {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//DO VALIDATION
	// check if user exists in group
	userIsMember := pipe.HExists(groupModel.GROUP_MEMBERS(groupID), userID)

	// get geo locations while in the round trip
	geoLocations := pipe.SMembers(groupModel.GROUP_GEO(groupID))

	_, err_pipe1 := pipe.Exec()

	if err_pipe1 != nil {
		return []string{}, errors.New("Pipe error.")
	}

	if !userIsMember.Val() {
		return []string{}, errors.New("You are not a member of this group.")
	}

	return geoLocations.Val(), nil
}
