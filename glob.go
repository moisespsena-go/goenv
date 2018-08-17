package goenv

import (
	"github.com/gobwas/glob"
	"github.com/gobwas/glob/syntax"
)

func globHasSpecial(s string) bool {
	for i := 0; i < len(s); i++ {
		if syntax.Special(s[i]) {
			return true
		}
	}
	return false
}

type singleGlob struct {
	s string
	g glob.Glob
}

func (s *singleGlob) Match(v string) bool {
	return s.g.Match(v)
}

func PathGlobCompile(s string) (g glob.Glob, err error) {
	if globHasSpecial(s) {
		return glob.Compile(s)
	} else {
		return &singleGlob{s, glob.MustCompile("**/" + s)}, nil
	}
}