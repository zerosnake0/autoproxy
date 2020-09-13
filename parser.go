package autoproxy

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

func ParseRule(rule string) (Rule, error) {
	if len(rule) == 0 {
		return nil, errors.New("empty rule")
	}
	if rule[0] == '|' {
		if len(rule) == 1 {
			return nil, fmt.Errorf("bad rule %q", rule)
		}
		if rule[1] == '|' {
			subRule := rule[2:]
			if subRule == "" {
				return nil, fmt.Errorf("empty domain rule %q", rule)
			}
			// ||example.com
			return domainRule{subRule}, nil
		} else {
			// len(rule) != 1, so non empty
			// |https://example.com
			return prefixRule{rule[1:]}, nil
		}
	}
	if rule[0] == '/' {
		if len(rule) <= 2 {
			return nil, fmt.Errorf("bad regex rule %q", rule)
		}
		if rule[len(rule)-1] != '/' {
			return nil, fmt.Errorf("bad regex rule %q", rule)
		}
		subRule := rule[1 : len(rule)-1]
		reg, err := regexp.Compile(subRule)
		if err != nil {
			return nil, fmt.Errorf("bad regex rule %q: %s", rule, err)
		}
		return regexRule{reg}, nil
	}
	return keywordRule{rule}, nil
}

func ParseRulesFromReader(r io.Reader) (rules, exceptions []Rule, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := bytes.TrimSpace(scanner.Bytes())
		if len(b) == 0 || b[0] == '!' || b[0] == '[' {
			continue
		}
		if bytes.HasPrefix(b, []byte{'@', '@'}) {
			// exception
			rule, err := ParseRule(string(b[2:]))
			if err != nil {
				return nil, nil, err
			}
			exceptions = append(exceptions, rule)
		} else {
			// match rules
			rule, err := ParseRule(string(b))
			if err != nil {
				return nil, nil, err
			}
			rules = append(rules, rule)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("unable to read proxy rules: %s", err)
	}
	return
}
