package groupModel

const (
	gr = "gr:" //group name maps to id

	grKeys = "grKeys"  //increments on each new group
	grHash = "grHash:" //group id maps to hash map for group information
	grMem  = "grMem:"  //set containing members for each group
	grInv  = "grInv"   //set containing members that currently have an invite to a group
)

func GROUP_NAME(s string) string {
	return gr + s
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
