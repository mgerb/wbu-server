package groupModel

const (
	grKeys = "grKeys" //increments on each new group

	groupID = "groupID" //hash - group name maps to id
	grHash  = "grHash:" //hash - group id maps to hash map for group information
	grMem   = "grMem:"  //hash - containing members for each group
	grMsg   = "grMsg:"  //list containing messages for each group id/username/message/timestamp
	grGeo   = "grGeo:"  //geo set for store user locations for each group
)

func GROUP_KEY_STORE() string {
	return grKeys
}

func GROUP_ID() string {
	return groupID
}

func GROUP_HASH(s string) string {
	return grHash + s
}

func GROUP_MEMBERS(s string) string {
	return grMem + s
}

func GROUP_MESSAGE(s string) string {
	return grMsg + s
}

func GROUP_GEO(s string) string {
	return grGeo + s
}

func GROUP_HASH_MAP(groupName string, owner string) map[string]string {
	return map[string]string{
		"groupName": groupName,
		"owner":     owner,
	}
}
