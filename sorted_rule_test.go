package autoproxy

import (
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/zerosnake0/autoproxy/internal/mock"
)

func TestSortedRule_Match(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockRule(ctrl)
	sr := sortedRule{rule: m}

	m.EXPECT().Match(gomock.Any()).Return(true)
	ok := sr.Match(nil)
	require.True(t, ok)
	require.Equal(t, uint64(1), sr.weight)

	m.EXPECT().Match(gomock.Any()).Return(false)
	ok = sr.Match(nil)
	require.False(t, ok)
	require.Equal(t, uint64(1), sr.weight)
}

func TestSortedRules_Append(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m1 := mock.NewMockRule(ctrl)
	m2 := mock.NewMockRule(ctrl)

	r1 := sortedRule{rule: m1, weight: 0}
	r2 := sortedRule{rule: m2, weight: 0}

	sorted := sortedAppendRules(nil, []Rule{m1, m2})
	require.Equal(t, 2, len(sorted))
	require.Equal(t, sorted[0], r1)
	require.Equal(t, sorted[1], r2)
}

func TestSortedRules_Sort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m1 := mock.NewMockRule(ctrl)
	m2 := mock.NewMockRule(ctrl)

	r1 := sortedRule{rule: m1, weight: 1}
	r2 := sortedRule{rule: m2, weight: 2}
	sorted := sortedRules{r1, r2}

	sort.Sort(sorted)
	require.Equal(t, sorted[0], r2)
	require.Equal(t, sorted[1], r1)
}
