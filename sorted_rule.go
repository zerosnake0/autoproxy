package autoproxy

import (
	"net/url"
	"sort"
	"sync/atomic"
)

type sortedRule struct {
	rule Rule

	weight uint64
}

var _ Rule = &sortedRule{}

func (rule *sortedRule) Match(u *url.URL) bool {
	res := rule.rule.Match(u)
	if res {
		atomic.AddUint64(&rule.weight, 1)
	}
	return res
}

type sortedRules []sortedRule

var _ sort.Interface = sortedRules{}

func (rules sortedRules) Len() int {
	return len(rules)
}

func (rules sortedRules) Less(i, j int) bool {
	return rules[i].weight > rules[j].weight
}

func (rules sortedRules) Swap(i, j int) {
	rules[i], rules[j] = rules[j], rules[i]
}

func sortedAppendRules(sorted sortedRules, rules []Rule) sortedRules {
	newRules := make(sortedRules, len(sorted), len(sorted)+len(rules))
	copy(newRules, sorted)
	for _, rule := range rules {
		newRules = append(newRules, sortedRule{
			rule: rule,
		})
	}
	return newRules
}
