package autoproxy

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexRule_Match(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := regexRule{}
		testRuleSet(t, r, nil, []string{
			"http://example.com",
		})
	})
	t.Run("test", func(t *testing.T) {
		reg, err := regexp.Compile("^https?:\\/\\/[^\\/]+example\\.com")
		require.NoError(t, err)
		require.NotNil(t, reg)
		r := regexRule{reg}
		testRuleSet(t, r, []string{
			"http://www.example.com",
			"https://www.example.com",
			"http://usr:pwd@www.example.com/search?a=b",
		}, []string{
			"http://www.example2.com",
		})
	})
}
