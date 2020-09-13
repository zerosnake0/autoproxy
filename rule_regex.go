package autoproxy

import (
	"net/url"
	"regexp"
)

type regexRule struct {
	rule *regexp.Regexp
}

var _ Rule = regexRule{}

func (r regexRule) Match(u *url.URL) bool {
	if r.rule == nil {
		return false
	}
	u = &url.URL{
		Scheme: u.Scheme,
		Opaque: u.Opaque,
		Host:   u.Host,
	}
	return r.rule.MatchString(u.String())
}
