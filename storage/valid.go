package storage

import (
	"errors"
	"unicode"
)

var (
	ErrInvalidTag  = errors.New("invalid tag; only letters, digits and underscores are allowed")
	ErrInvalidPath = errors.New("invalid path; control characters are not allowed")
)

func IsValidTag(tag string) bool {
	for _, chr := range tag {
		if !unicode.In(chr, unicode.Letter, unicode.Digit) && chr != '_' {
			return false
		}
	}
	return true
}

func IsValidPath(path string) bool {
	for _, chr := range path {
		if unicode.IsControl(chr) {
			return false
		}
	}
	return true
}
