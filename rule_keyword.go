package autoproxy

import (
	"net/url"
	"strings"
)

type keywordRule struct {
	rule string
}

var _ Rule = keywordRule{}

func (r keywordRule) Match(u *url.URL) bool {
	if r.rule == "" {
		return false
	}
	if u.Scheme != "http" {
		return false
	}
	u = removeUser(u)
	return strings.Contains(u.String(), r.rule)
}
