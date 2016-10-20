package groupOperations

import (
	"errors"
	"strconv"

	"../../db"
	"../../model/groupModel"
	redis "gopkg.in/redis.v4"
)

func StoreGeoLocation(groupID string, longitude string, latitude string, userID string, username string) error {

	//DO VALIDATION
	//check if user exists in group before storing message
	if !UserIsMember(userID, groupID) {
		return errors.New("User is not in group.")
	}

	long, err1 := strconv.ParseFloat(longitude, 64)
	lat, err2 := strconv.ParseFloat(latitude, 64)

	if err1 != nil || err2 != nil {
		return errors.New("Invalid longitude/latitude")
	}

	geoLocation := &redis.GeoLocation{}

	geoLocation.Name = userID + "/" + username
	geoLocation.Longitude = long
	geoLocation.Latitude = lat

	_, returnError := db.Client.GeoAdd(groupModel.GROUP_GEO(groupID), geoLocation).Result()

	return returnError
}