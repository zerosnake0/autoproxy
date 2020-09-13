package autoproxy

import (
	"io"
	"net/url"
	"sort"
	"sync"
	"time"
)

type AutoProxy struct {
	mu sync.RWMutex

	sortPeriod time.Duration
	lastSort   time.Time
	ch         chan struct{}

	rules      sortedRules
	exceptions sortedRules
}

type Option struct {
	SortPeriod time.Duration
}

func New(opt *Option) *AutoProxy {
	p := &AutoProxy{
		ch: make(chan struct{}, 1),
	}
	if opt != nil {
		p.sortPeriod = opt.SortPeriod
	}
	p.readyToSort()
	return p
}

func (p *AutoProxy) Read(r io.Reader) error {
	rules, exceptions, err := ParseRulesFromReader(r)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = sortedAppendRules(p.rules, rules)
	p.exceptions = sortedAppendRules(p.exceptions, exceptions)
	return nil
}

func (p *AutoProxy) readyToSort() {
	p.ch <- struct{}{}
}

func (p *AutoProxy) sort() {
	select {
	case <-p.ch:
	default:
		return
	}
	go func() {
		defer p.readyToSort()
		p.mu.Lock()
		defer p.mu.Unlock()
		elapsed := time.Now().Sub(p.lastSort)
		if elapsed > p.sortPeriod {
			sort.Reverse(p.rules)
			sort.Reverse(p.exceptions)
		}
	}()
}

func (p *AutoProxy) Match(u *url.URL) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.sortPeriod > 0 {
		elapsed := time.Now().Sub(p.lastSort)
		if elapsed > p.sortPeriod {
			p.sort()
		}
	}
	for _, rule := range p.exceptions {
		if rule.Match(u) {
			return false
		}
	}
	for _, rule := range p.rules {
		if rule.Match(u) {
			return true
		}
	}
	return false
}
