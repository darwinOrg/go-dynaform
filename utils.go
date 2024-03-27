package dynaform

import "strings"

func FirstCharToUpper(s string) string {
	if s == "" {
		return ""
	}

	if len(s) == 1 {
		return strings.ToUpper(s)
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}
