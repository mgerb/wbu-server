package model

import (
	"strconv"
)

//schema
const (
	u = "u:" //username maps to user id

	u_keys   = "u_keys"    //increments on each new user
	u_hash   = "u_hash:"   //hash map for users = id maps to key values
	u_groups = "u_groups:" //sorted set for groups that user is in
)

func USER_GROUPS(s string) string {
	return u_groups + s
}

func USER_NAME(s string) string {
	return u + s
}

func USER_HASH(i int) string {
	return u_hash + strconv.Itoa(i)
}

func USER_KEY_STORE() string {
	return u_keys
}
