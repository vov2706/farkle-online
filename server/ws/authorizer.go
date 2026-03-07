package ws

import (
	"fmt"
	"regexp"
	"strings"
)

type AuthorizeFunc func(user any, params map[string]string) (bool, map[string]any)

type Rule struct {
	Pattern *regexp.Regexp
	Params  []string
	Check   AuthorizeFunc
}

type Authorizer struct {
	rules []Rule
}

func NewAuthorizer() *Authorizer { return &Authorizer{} }

func (a *Authorizer) Channel(pattern string, paramNames []string, check AuthorizeFunc) {
	regexStr := "^" + pattern + "$"

	for _, p := range paramNames {
		regexStr = strings.ReplaceAll(regexStr, "{"+p+"}", "([^:]+)")
	}

	a.rules = append(a.rules, Rule{
		Pattern: regexp.MustCompile(regexStr),
		Params:  paramNames,
		Check:   check,
	})
}

func (a *Authorizer) Authorize(user any, channel string) (bool, map[string]any, error) {
	for _, r := range a.rules {
		m := r.Pattern.FindStringSubmatch(channel)

		if m == nil {
			continue
		}

		if len(m)-1 != len(r.Params) {
			return false, nil, fmt.Errorf("authorizer param mismatch for channel %s", channel)
		}
		params := map[string]string{}
		for i, name := range r.Params {
			params[name] = m[i+1]
		}

		ok, meta := r.Check(user, params)
		return ok, meta, nil
	}
	return false, nil, nil // no rule => forbidden
}
