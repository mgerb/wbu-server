package geoRoutes

import (
	"../../operations/geoOperations"
	"../../utils/response"
	"github.com/labstack/echo"
)

// StoreGeoLocation -
func StoreGeoLocation(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	groupID := ctx.FormValue("groupID")
	latitude := ctx.FormValue("latitude")
	longitude := ctx.FormValue("longitude")

	err := geoOperations.StoreGeoLocation(userID, groupID, latitude, longitude)

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
