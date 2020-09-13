package autoproxy

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

func ParseRule(rule string) (Rule, error) {
	if len(rule) == 0 {
		return nil, errors.New("empty rule")
	}
	if rule[0] == '|' {
		if len(rule) == 1 {
			return nil, errors.New("bad rule")
		}
		if rule[1] == '|' {
			subRule := rule[2:]
			if subRule == "" {
				return nil, errors.New("empty domain rule")
			}
			// ||example.com
			return domainRule{subRule}, nil
		} else {
			// len(rule) != 1, so non empty
			// |https://example.com
			return prefixRule{rule[1:]}, nil
		}
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
		return nil, nil, fmt.Errorf("unable to read proxy rules: %w", err)
	}
	return
}
