package autoproxy

import "testing"

func TestPrefixRule_Match(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := prefixRule{}
		testRuleSet(t, r, nil, []string{
			"http://example.com",
		})
	})
	t.Run("test", func(t *testing.T) {
		r := prefixRule{"http://example.com"}
		testRuleSet(t, r, []string{
			"http://example.com",
			"http://usr:pwd@example.com",
		}, []string{
			"https://example.com",
		})
	})
}
