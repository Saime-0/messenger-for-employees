package utils

import (
	"strings"
)

func ClearMessageBodyOfExtraCharacters(body *string) *string {
	if body == nil {
		return nil
	}
	var pureBody = ""
	var lines = strings.Split(*body, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			pureBody += line
			if i+1 != len(lines) {
				pureBody += "\n"
			}
		}
	}
	return &pureBody
}
