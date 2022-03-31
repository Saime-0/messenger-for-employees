package validator

import (
	"github.com/saime-0/messenger-for-employee/internal/rules"
	"regexp"
	"strings"
)

var (
	expEmail      = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	expRoomName   = regexp.MustCompile(`^([a-z]{2,16}|[а-я]{2,16})+$`)
	expSessionKey = regexp.MustCompile(`^[a-zA-Z0-9\-=]{20}$`)
	expPartOfName = regexp.MustCompile(`^([a-z]{2,16}|[а-я]{2,16})+$`)
)

func ValidateEmployeeFullName(fullName string) (valid bool) {
	for i, partOfName := range strings.Split(strings.ToLower(fullName), " ") {
		if !expPartOfName.MatchString(partOfName) || i > 1 {
			return false
		}
	}
	return true
}

func ValidatePartOfName(part string) (valid bool) {
	return expPartOfName.MatchString(part)
}

func ValidateRoomName(name string) (valid bool) {
	return expRoomName.MatchString(name)
}

func ValidatePassword(password string) (valid bool) { // todo
	return len(password) <= rules.MaxPasswordLength && len(password) >= rules.MinPasswordLength
}

func ValidateEmail(email string) (valid bool) { // todo
	return expEmail.MatchString(email)
}

func ValidateOffset(offset int) (valid bool) {
	return offset >= 0
}

func ValidateLimit(limit int) (valid bool) {
	return limit >= 1 && limit <= 20
}

func ValidateID(id int) (valid bool) {
	return id > 0
}

func ValidateIDs(ids []int) (valid bool) {
	for _, id := range ids {
		if !ValidateID(id) {
			return false
		}
	}
	return true
}

func ValidateSessionKey(sessionKey string) (valid bool) {
	return expSessionKey.MatchString(sessionKey)
}
