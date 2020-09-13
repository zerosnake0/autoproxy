package autoproxy

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestAutoProxy_Read(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		p := New(nil)
		err := p.Read(strings.NewReader("@@"))
		require.Error(t, err)
	})
	t.Run("concurrent", func(t *testing.T) {
		p := New(nil)
		count := 10
		var eg errgroup.Group
		for i := 0; i < count; i++ {
			raw := fmt.Sprintf("example%d.com", i)
			eg.Go(func() error {
				return p.Read(strings.NewReader(raw))
			})
		}
		err := eg.Wait()
		require.NoError(t, err)
		for i := 0; i < count; i++ {
			raw := fmt.Sprintf("http://example%d.com", i)
			ok, err := MatchRule(p, raw)
			require.NoError(t, err)
			require.True(t, ok)
		}
	})
	t.Run("file", func(t *testing.T) {
		must := require.New(t)
		fn := os.Getenv("AUTOPROXY_TEST_FILE")
		if fn == "" {
			t.SkipNow()
		}
		f, err := os.Open(fn)
		must.NoError(err)
		defer f.Close()
		p := New(nil)
		err = p.Read(f)
		must.NoError(err)
	})
}

func TestAutoProxy_Match(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		must := require.New(t)
		p := New(nil)
		err := p.Read(strings.NewReader(`[AutoProxy]
		!comment
		example.com
		|http://example2.com
		||example3.com
		@@||example4.com`))
		must.NoError(err)
		testRuleSet(t, p, []string{
			"http://example.com",
		}, []string{
			"https://example.com",
			"http://subdomain.example4.com",
			"http://example5.com",
		})
	})
	t.Run("sort", func(t *testing.T) {
		must := require.New(t)
		p := New(&Option{
			SortPeriod: time.Second,
		})
		err := p.Read(strings.NewReader(`
		example.com
		example2.com`))
		must.NoError(err)
		var eg errgroup.Group
		for i, _raw := range []string{
			"http://example.com",
			"http://example2.com",
		} {
			raw := _raw
			for j := 0; j <= i; j++ {
				eg.Go(func() error {
					ok, err := MatchRule(p, raw)
					if err != nil {
						return err
					}
					if !ok {
						return errors.New("should not reach here")
					}
					return nil
				})
			}
		}
		err = eg.Wait()
		must.NoError(err)
	})
}
