package regex

//regular expressions
const (
	USERNAME   = `^[a-zA-Z1-9_\s]{2,20}$`     //only allow characters and underscores
	GROUP_NAME = `^\S[a-zA-Z1-9_\s]{0,20}\S$` //allow characters, underscores, and spaces in middle
	PASSWORD   = `^[a-zA-Z1-9!@#$%^&*_.\-]{5,20}$`
	EMAIL      = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)
