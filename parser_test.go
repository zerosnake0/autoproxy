package autoproxy

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRule(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r, err := ParseRule("")
		require.Error(t, err)
		require.Nil(t, r)
	})
	t.Run("bad rule |", func(t *testing.T) {
		r, err := ParseRule("|")
		require.Error(t, err)
		require.Nil(t, r)
	})
	t.Run("empty domain", func(t *testing.T) {
		r, err := ParseRule("||")
		require.Error(t, err)
		require.Nil(t, r)
	})
	t.Run("domain", func(t *testing.T) {
		r, err := ParseRule("||example.com")
		require.NoError(t, err)
		require.Equal(t, domainRule{"example.com"}, r)
	})
	t.Run("prefix", func(t *testing.T) {
		r, err := ParseRule("|example.com")
		require.NoError(t, err)
		require.Equal(t, prefixRule{"example.com"}, r)
	})
	t.Run("keyword", func(t *testing.T) {
		kw := "example.com"
		r, err := ParseRule(kw)
		require.NoError(t, err)
		require.Equal(t, keywordRule{kw}, r)
	})
	t.Run("regex", func(t *testing.T) {
		t.Run("bad", func(t *testing.T) {
			r, err := ParseRule("//")
			require.Error(t, err)
			require.Nil(t, r)
		})
		t.Run("bad", func(t *testing.T) {
			r, err := ParseRule("/^a")
			require.Error(t, err)
			require.Nil(t, r)
		})
		t.Run("bad", func(t *testing.T) {
			r, err := ParseRule("/[^]/")
			require.Error(t, err)
			require.Nil(t, r)
		})
		t.Run("good", func(t *testing.T) {
			raw := "^abc$"
			r, err := ParseRule("/" + raw + "/")
			require.NoError(t, err)
			reg, err := regexp.Compile(raw)
			require.NoError(t, err)
			require.Equal(t, regexRule{reg}, r)
		})
	})
}

type badReader struct{}

var _ io.Reader = badReader{}

func (badReader) Read([]byte) (int, error) {
	return 0, errors.New("test")
}

func TestParseRulesFromReader(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("reader", func(t *testing.T) {
			_, _, err := ParseRulesFromReader(badReader{})
			require.Error(t, err)
		})
		t.Run("bad rule", func(t *testing.T) {
			_, _, err := ParseRulesFromReader(strings.NewReader("|"))
			require.Error(t, err)
		})
		t.Run("bad exception", func(t *testing.T) {
			_, _, err := ParseRulesFromReader(strings.NewReader("@@"))
			require.Error(t, err)
		})
	})
	t.Run("test", func(t *testing.T) {
		r := strings.NewReader(`
			[AutoProxy]
			!comment
			example.com
			@@||example.com
		`)
		rules, exceptions, err := ParseRulesFromReader(r)
		require.NoError(t, err)

		require.Equal(t, 1, len(rules))
		require.Equal(t, keywordRule{"example.com"}, rules[0])

		require.Equal(t, 1, len(exceptions))
		require.Equal(t, domainRule{"example.com"}, exceptions[0])
	})
}
