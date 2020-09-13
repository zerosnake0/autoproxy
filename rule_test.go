package autoproxy

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/zerosnake0/autoproxy/internal/mock"
)

func testRuleSet(t *testing.T, rule Rule, matches, nonMatches []string) {
	f := func(exp bool, arr []string) {
		for _, url := range arr {
			got, err := MatchRule(rule, url)
			require.NoError(t, err)
			require.Equal(t, exp, got, "%q expect %t but got %t", url, exp, got)
		}
	}
	f(true, matches)
	f(false, nonMatches)
}

func TestMatchRule(t *testing.T) {
	t.Run("wrong url", func(t *testing.T) {
		_, err := MatchRule(nil, "http://example.com:badport")
		require.Error(t, err)
	})
	t.Run("test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock.NewMockRule(ctrl)

		raw := "http://example.com"

		m.EXPECT().Match(gomock.Any()).Return(true)
		ok, err := MatchRule(m, raw)
		require.NoError(t, err)
		require.True(t, ok)

		m.EXPECT().Match(gomock.Any()).Return(false)
		ok, err = MatchRule(m, raw)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestRemoveUser(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		exp := &url.URL{}
		got := removeUser(exp)
		require.Equal(t, exp, got)
	})
	t.Run("non nil", func(t *testing.T) {
		input := &url.URL{User: &url.Userinfo{}}
		got := removeUser(input)
		require.NotNil(t, input.User)
		require.Nil(t, got.User)
	})

}
