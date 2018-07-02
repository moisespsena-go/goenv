package cmd

import "strings"

func pad(v string, s ...int) string {
	l := 10
	if len(s) > 0 {
		l = s[0]
	}
	if len(v) >= l {
		return v
	}
	return v + strings.Repeat(" ", l-len(v))
}
