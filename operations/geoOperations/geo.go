package geoOperations

import (
	"errors"
	"log"
	"strconv"

	"github.com/mgerb/wbu-server/db"
	"github.com/mgerb/wbu-server/model"
)

// StoreGeoLocation -
func StoreGeoLocation(userID string, groupID string, latitude string, longitude string, waypoint bool) error {

	// convert coordinates for validation
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

	//check if user exists in group before storing location
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	defer tx.Commit()

	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE userID = ? AND groupID = ?);`, userID, groupID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	if !userExistsInGroup {
		return errors.New("user not in group")
	}

	// insert/update user location
	_, err = tx.Exec(`UPDATE "GeoLocation" SET latitude = ?, longitude = ?, timestamp = CURRENT_TIMESTAMP, waypoint = ? WHERE userID = ? AND groupID = ? AND waypoint = ?;
					INSERT INTO "GeoLocation" (userID, groupID, latitude, longitude, waypoint) SELECT ?, ?, ?, ?, ? WHERE changes() = 0;`,
		latitude, longitude, waypoint, userID, groupID, waypoint, userID, groupID, latitude, longitude, waypoint)

	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

// GetGeoLocations -
func GetGeoLocations(userID string, groupID string) ([]*model.GeoLocation, error) {

	// start SQL transaction
	tx, err := db.SQL.Begin()
	if err != nil {
		log.Println(err)
		return []*model.GeoLocation{}, errors.New("database error")
	}

	defer tx.Commit()

	// check if user exists in group
	var userExistsInGroup bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM "UserGroup" WHERE "userID" = ? AND "groupID" = ?);`, userID, groupID).Scan(&userExistsInGroup)

	if err != nil {
		log.Println(err)
		return []*model.GeoLocation{}, errors.New("database error")
	}

	if !userExistsInGroup {
		return []*model.GeoLocation{}, errors.New("user not in group")
	}

	// get GeoLocation
	rows, err := tx.Query(`SELECT u.id, u.firstName, u.lastName, u.email, gl.id, gl.groupID, gl.latitude, gl.longitude, strftime('%s', gl.timestamp), gl.waypoint
							FROM "GeoLocation" AS gl INNER JOIN "User" AS u ON gl.userID = u.id
							WHERE gl.groupID = ?;`, groupID)

	if err != nil {
		log.Println(err)
		return []*model.GeoLocation{}, errors.New("database error")
	}

	defer rows.Close()

	geoList := []*model.GeoLocation{}

	for rows.Next() {
		newGeo := &model.GeoLocation{}
		err := rows.Scan(&newGeo.UserID, &newGeo.FirstName, &newGeo.LastName, &newGeo.Email, &newGeo.ID,
			&newGeo.GroupID, &newGeo.Latitude, &newGeo.Longitude, &newGeo.Timestamp, &newGeo.Waypoint)

		if err != nil {
			log.Println(err)
			return []*model.GeoLocation{}, errors.New("database error")
		}

		geoList = append(geoList, newGeo)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return []*model.GeoLocation{}, errors.New("database error")
	}

	// return geo model and nil
	return geoList, nil
}
