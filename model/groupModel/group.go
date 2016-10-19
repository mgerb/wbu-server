package groupModel

const (
	grKeys = "grKeys"  //increments on each new group
	
	grID   = "grID:"   //group name maps to id
	grHash = "grHash:" //group id maps to hash map for group information
	grMem  = "grMem:"  //set containing members for each group
	grInv  = "grInv"   //set containing members that currently have an invite to a group
)

func GROUP_ID(s string) string {
	return grID + s
}

func GROUP_HASH(s string) string {
	return grHash + s
}

func GROUP_MEMBERS(s string) string {
	return grMem + s
}

func GROUP_KEY_STORE() string {
	return grKeys
}
