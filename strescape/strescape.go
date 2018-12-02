// Package strescape is utility tool to escape strings.
package strescape

import "strings"

// SingleQuote convert all "'" -> "\\'".
func SingleQuote(str string) string {
	return strings.Replace(str, "'", "\\'", -1)
}
