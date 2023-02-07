package pac

import "strings"

type ruleErr []string

// append error message for proxy rule
func rerrAppend(prev ruleErr, msg string) ruleErr {
	if prev == nil {
		prev = (ruleErr)(make([]string, 0, 4))
	}
	return append(prev, msg)
}

// error report
func (re ruleErr) Error() string {
	return "Rule error:\n    " + strings.Join(re, "\n    ")
}

// exprot error
func (re ruleErr) Export() error {
	if re == nil {
		return nil
	}
	return re
}
