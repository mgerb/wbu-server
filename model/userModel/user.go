package userModel

//schema
const (
	usr = "usr:" //username maps to user id

	uKeys   = "uKeys"    //increments on each new user
	uHash   = "uHash:"   //hash map for users = id maps to key values
	uGroups = "uGroups:" //sorted set for groups that user is in
)

func USER_GROUPS(s string) string {
	return uGroups + s
}

func USER_NAME(s string) string {
	return usr + s
}

func USER_HASH(s string) string {
	return uHash + s
}

func USER_KEY_STORE() string {
	return uKeys
}
