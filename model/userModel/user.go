package userModel

//schema
const (
	uKeys   = "uKeys"    //increments on each new user

	uID     = "uID:"     //username maps to user id
	uHash   = "uHash:"   //hash map for users = id maps to key values
	uGroups = "uGroups:" //sorted set for groups that user is in
)

func USER_ID(s string) string {
	return uID + s
}

func USER_KEY_STORE() string {
	return uKeys
}

func USER_HASH(s string) string {
	return uHash + s
}

func USER_GROUPS(s string) string {
	return uGroups + s
}