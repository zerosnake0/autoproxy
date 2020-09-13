package autoproxy

import (
	"testing"
)

func TestDomainRule_Match(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := domainRule{}
		testRuleSet(t, r, nil, []string{
			"http://example.com",
		})
	})
	t.Run("test", func(t *testing.T) {
		r := domainRule{"example.com"}
		testRuleSet(t, r, []string{
			"http://example.com",
			"https://example.com",
			"https://subdomain.example.com",
			"http://usr:pwd@subdomain.example.com",
		}, []string{
			"http://example2.com/search?example.com",
		})
	})
}
