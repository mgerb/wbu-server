package userModel

//schema
const (
	uKeys = "uKeys" //increments on each new user

	uID     = "uID:"     //username maps to user id
	uHash   = "uHash:"   //hash map for users = id maps to key values
	uGrps   = "uGrps:"   //sorted set for groups that user is in
	uGrpInv = "uGrpInv:" //set of pending group invites
)

func USER_KEY_STORE() string {
	return uKeys
}

func USER_ID(s string) string {
	return uID + s
}

func USER_HASH(s string) string {
	return uHash + s
}

func USER_GROUPS(s string) string {
	return uGrps + s
}

func USER_GROUP_INVITES(s string) string {
	return uGrpInv + s
}

func USER_HASH_MAP(userName string, password string, adminGroupCount string) map[string]string {
	return map[string]string{
		"userName":        userName,
		"password":        password,
		"email":           "string",
		"adminGroupCount": adminGroupCount,
		"fcmToken":        "string",
		"facebookToken":   "string",
	}
}
