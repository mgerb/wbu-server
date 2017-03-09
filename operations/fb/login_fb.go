//facebook operations
package fb

import (
	fb "github.com/huandu/facebook"
)

func Me(accessToken string) (map[string]interface{}, error) {
	return fb.Get("/me", fb.Params{
		"fields":       "id,first_name,last_name,email",
		"access_token": accessToken,
	})
}
