package kecpvalidate

import "unicode"

func IsValidUserName(s string) bool {
	if len(s) > 16 {
		return false
	}
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
