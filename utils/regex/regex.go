package regex

//regular expressions
const (
	USERNAME   = `^[a-zA-Z0-9_\s]{2,20}$`      //only allow characters and underscores
	GROUP_NAME = `^\S[a-zA-Z0-9'_\s]{0,20}\S$` //allow characters, underscores, and spaces in middle
	PASSWORD   = `^[a-zA-Z0-9!@#$%^&*_.\-]{5,20}$`
	EMAIL      = `^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)
