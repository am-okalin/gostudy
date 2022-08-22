package strings

import (
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	s1 := "aBCDEcf"
	s2 := strings.ToLower(s1)

	i1 := strings.Index(s2, "c")
	t.Log(i1)

	i2 := strings.LastIndex(s2, "c")
	t.Log(i2)

	before, after, found := strings.Cut(s2, "c")
	t.Log(before, after, found)
}
