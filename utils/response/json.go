package response

//error types
const (
	SUCCESS                = 1
	INTERNAL_ERROR         = 2
	INVALID_AUTHENTICATION = 3
	//...
	//...
)

//Json - return json string for response
func Json(message string, statusCode int) map[string]interface{} {

	return map[string]interface{}{
		"status":  statusCode,
		"message": string(message),
	}
}
