package response

//error types
const (
	INTERNAL_ERROR = 2
	SUCCESS        = 1
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
