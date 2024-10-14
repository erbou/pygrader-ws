package uti

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	ErrCNameInvalid ErrorCode = 1100 + iota
	ErrEmailInvalid
)

var validEmailRFC5322 = regexp.MustCompile(`^(?i)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)

var tknmEmail = regexp.MustCompile(`\+.*$|\.+|^[^:]*:`)

var validCName = regexp.MustCompile(`^(?i)[A-Z0-9](?:[._A-Z0-9 -]*[A-Z0-9])?$`)

var wSpace = regexp.MustCompile(`\s+`)

func CanonizeEmail(s string) (string, error) {
	if !validEmailRFC5322.MatchString(s) {
		return ``, Errorf(ErrEmailInvalid, `Invalid email '%v'`, s)
	}

	s = strings.ToLower(s)
	s = strings.ToValidUTF8(s, `_`)
	email := strings.Split(s, `@`)
	s = tknmEmail.ReplaceAllString(email[0], ``) + `@` + email[1]

	return s, nil
}

func CanonizeName(s string) (string, error) {
	tc := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(tc, s)

	if !validCName.MatchString(s) {
		return ``, Errorf(ErrCNameInvalid, `Invalid name '%v'`, s)
	}

	s = strings.ToLower(s)
	s = strings.ToValidUTF8(s, `_`)
	s = wSpace.ReplaceAllString(s, ` `)

	return s, nil
}
