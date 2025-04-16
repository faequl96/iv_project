package utils

import "strings"

func SplitName(fullname string) (string, string) {
	parts := strings.Fields(fullname)
	if len(parts) == 0 {
		return "", ""
	}

	firstname := parts[0]
	lastname := ""
	if len(parts) > 1 {
		lastname = strings.Join(parts[1:], " ")
	}

	return firstname, lastname
}
