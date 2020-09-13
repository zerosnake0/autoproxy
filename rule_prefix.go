package autoproxy

import (
	"net/url"
	"strings"
)

type prefixRule struct {
	rule string
}

var _ Rule = prefixRule{}

func (r prefixRule) Match(u *url.URL) bool {
	if r.rule == "" {
		return false
	}
	u = removeUser(u)
	return strings.HasPrefix(u.String(), r.rule)
}
