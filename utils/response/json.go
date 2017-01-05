package response

//error types
const (
	SUCCESS                int    = 1
	INTERNAL_ERROR         int    = 2
	INVALID_AUTHENTICATION int    = 3
	JSON_HEADER            string = "application/json; charset=utf-8"
)

//Json - return json string for response
func Json(message string, statusCode int) map[string]interface{} {

	return map[string]interface{}{
		"status":  statusCode,
		"message": string(message),
	}
}
