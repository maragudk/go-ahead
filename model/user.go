package model

import (
	"regexp"
	"time"
)

type ID = string
type Email = string

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidateEmail according to the emailRegexp.
func ValidateEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

// User with a Name, Email, and Password.
type User struct {
	ID       ID
	Name     string
	Email    Email
	Password string `json:"-"`
	Created  time.Time
	Updated  time.Time
}
