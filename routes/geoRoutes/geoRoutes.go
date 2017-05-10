package geoRoutes

import (
	"github.com/labstack/echo"
	"github.com/mgerb/wbu-server/operations/geoOperations"
	"github.com/mgerb/wbu-server/utils/response"
)

// StoreGeoLocation -
func StoreGeoLocation(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	latitude := ctx.FormValue("latitude")
	longitude := ctx.FormValue("longitude")
	waypoint := ctx.FormValue("waypoint") == "true"

	err := geoOperations.StoreGeoLocation(userID, groupID, latitude, longitude, waypoint)

	switch err {
	case nil:
		return ctx.JSON(200, response.Json("Geo location updated.", response.SUCCESS))
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}

// GetGeoLocations -
func GetGeoLocations(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.Param("groupID")

	geoLocations, err := geoOperations.GetGeoLocations(userID, groupID)

	switch err {
	case nil:
		return ctx.JSON(200, geoLocations)
	default:
		return ctx.JSON(500, response.Json(err.Error(), response.INTERNAL_ERROR))
	}
}
