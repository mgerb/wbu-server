package groupOperations

import (
	"errors"
	"strconv"

	"../../db"
)

// StoreGeoLocation -
func StoreGeoLocation(userID string, groupID string, latitude string, longitude string) error {

	pipe := db.Client.Pipeline()
	defer pipe.Close()

	//DO VALIDATION
	//check if user exists in group before storing message

	long, err1 := strconv.ParseFloat(longitude, 64)
	lat, err2 := strconv.ParseFloat(latitude, 64)

	if err1 != nil || err2 != nil {
		return errors.New("invalid coordinates")
	}

	//validate coordinates
	if long < -180 ||
		long > 180 ||
		lat < -85.05112878 ||
		lat > 85.05112878 {
		return errors.New("invalid coordinates")
	}

	return nil
}

/*
// GetGeoLocations -
func GetGeoLocations(groupID string, userID string) ([]string, error) {


	// get geo locations from geo table

	// return geo model and nil
	return geoLocations.Val(), nil
}

*/
