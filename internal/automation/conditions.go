package automation

import (
	"regexp"
	"strings"
)

// checkStringMatch evaluates basic string operators.
func checkStringMatch(target, operator, value string) bool {
	switch operator {
	case "EQUALS":
		return target == value
	case "NOT_EQUALS":
		return target != value
	case "CONTAINS":
		return strings.Contains(target, value)
	default:
		return false
	}
}

// checkRegexMatch compiles and executes regular expressions.
// NOTE: In production, heavily used regex patterns should be pre-compiled into a cache map to save CPU cycles.
func checkRegexMatch(target, pattern string) bool {
	matched, err := regexp.MatchString(pattern, target)
	if err != nil {
		return false
	}
	return matched
}