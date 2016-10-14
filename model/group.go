package model

//schema
const (
	gr_keys    = "gr_keys"     //increments on each new group
	gr         = "gr:"    //username maps to group id
	gr_hash    = "gr_hash:"    //group id maps to hash map for group information
	gr_members = "gr_members:" //set containing members for each group
	gr_invites = "gr_invites"  //set containing members that currently have an invite to a group
)
