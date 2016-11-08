package userModel

//schema
const (
	uKeys = "uKeys" //increments on each new user

	userID  = "userID"   //email hash maps to user id
	uHash   = "uHash:"   //hash map for users = id maps to key values
	uGrps   = "uGrps:"   //hash for user groups
	uGrpInv = "uGrpInv:" //hash for user group invites
)

func USER_KEY_STORE() string {
	return uKeys
}

func USER_ID() string {
	return userID
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

func USER_HASH_MAP(userName string, password string, adminGroupCount string, name string, email string, facebookID string) map[string]string {
	return map[string]string{
		"userName":        userName,
		"password":        password,
		"adminGroupCount": adminGroupCount,
		"fcmToken":        "",
		"name":            name,
		"email":           email,
		"facebookID":      facebookID,
	}
}
