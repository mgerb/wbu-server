package regex

//regular expressions
const (
	USERNAME   = "^[a-zA-Z1-9_]{2,15}$"
	GROUP_NAME = "^\\S[a-zA-Z1-9_\\s]{2,20}\\S$"
	PASSWORD   = "^[a-zA-Z1-9!@#$%^&*_.\\-]{5,20}$"
)
