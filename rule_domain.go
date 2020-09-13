package autoproxy

import (
	"net/url"
	"strings"
)

type domainRule struct {
	rule string
}

var _ Rule = domainRule{}

func (r domainRule) Match(u *url.URL) bool {
	host := u.Hostname()
	if !strings.HasSuffix(host, r.rule) {
		return false
	}
	idx := len(host) - len(r.rule)
	if idx == 0 {
		return true
	}
	return host[idx-1] == '.'
}
