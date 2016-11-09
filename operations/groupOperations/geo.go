package groupOperations

import (
	"../../db"
	"../../model/groupModel"
	"errors"
	redis "gopkg.in/redis.v5"
	"strconv"
	"time"
)

func StoreGeoLocation(groupID string, longitude string, latitude string, userID string, userName string) error {

	//DO VALIDATION
	//check if user exists in group before storing message
	userIsMember := db.Client.HExists(groupModel.GROUP_MEMBERS(groupID), userID).Val()

	if !userIsMember {
		return errors.New("You are not a member of this group")
	}

	long, err1 := strconv.ParseFloat(longitude, 64)
	lat, err2 := strconv.ParseFloat(latitude, 64)

	if err1 != nil || err2 != nil {
		return errors.New("Invalid longitude/latitude")
	}

	//validate coordinates
	if long < -180 ||
		long > 180 ||
		lat < -85.05112878 ||
		lat > 85.05112878 {
		return errors.New("Invalid coordinates.")
	}

	geoLocation := &redis.GeoLocation{}

	geoLocation.Name = userID + "/" + userName + "/" + strconv.FormatInt(time.Now().Unix(), 10)
	geoLocation.Longitude = long
	geoLocation.Latitude = lat

	_, returnError := db.Client.GeoAdd(groupModel.GROUP_GEO(groupID), geoLocation).Result()

	return returnError
}
