package autoproxy

import (
	"net/url"
)

type Rule interface {
	Match(u *url.URL) bool
}

func MatchRule(rule Rule, raw string) (bool, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return false, err
	}
	return rule.Match(u), nil
}

func removeUser(u *url.URL) *url.URL {
	if u.User == nil {
		return u
	}
	cpy := *u
	cpy.User = nil
	return &cpy
}
