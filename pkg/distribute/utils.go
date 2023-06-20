package distribute

import (
	"regexp"
)

const email_regex_str = `^\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b$`

var email_regex = regexp.MustCompile(email_regex_str)

func IsValidEmail(email string) bool {
	return email_regex.MatchString(email)
}
