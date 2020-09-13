package autoproxy

import (
	"testing"
)

func TestKeywordRule_Match(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := keywordRule{}
		testRuleSet(t, r, nil, []string{
			"http://example.com",
		})
	})
	t.Run("test", func(t *testing.T) {
		r := keywordRule{"example.com"}
		testRuleSet(t, r, []string{
			"http://example.com",
			"http://www.example2.com/search?q=example.com",
		}, []string{
			"https://example.com",
		})
	})
}
