package response

import (
	"encoding/json"
)

//error types
const (
	INTERNAL_ERROR = 1
	SUCCESS        = 2
	//...
	//...
)

//Json - return json string for response
func Json(message string, statusCode int) string {

	jsonMap := map[string]interface{}{
		"status":  statusCode,
		"message": string(message),
	}

	response, err := json.Marshal(jsonMap)

	if err != nil {
		return `{"error": "error"}`
	}

	return string(response)
}
