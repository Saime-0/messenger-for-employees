package validator

import (
	"github.com/saime-0/messenger-for-employee/internal/rules"
	"regexp"
	"strconv"
	"strings"
)

var (
	expEmail      = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	expRoomName   = regexp.MustCompile(`^.{1,64}$`)
	expTagName    = regexp.MustCompile(`^.{1,32}$`)
	expSessionKey = regexp.MustCompile(`^[a-zA-Z0-9\-=]{20}$`)
	expPartOfName = regexp.MustCompile(`^([a-z]{2,16}|[а-я]{2,16})+$`)
)

func ValidateEmployeeFullName(fullName string) (valid bool) {
	var names = strings.Split(strings.ToLower(fullName), " ")
	if len(names) != 2 {
		return false
	}
	for i, partOfName := range names {
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

func ValidateTagName(name string) (valid bool) {
	return expTagName.MatchString(name)
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

func ValidateMessageBody(body *string) (valid bool, reason string) {
	switch {
	case body == nil:
		return false, "ошибка произошла там где ее никто не ждал"
	case len(*body) < 1:
		return false, "нельзя отправлять пустое сообщение"
	case len(*body) > rules.MaxMessageBodyLen:
		return false, "превышен лимит символов, макс. длина:" + strconv.Itoa(rules.MaxMessageBodyLen)
	default:
		return true, ""
	}
}
