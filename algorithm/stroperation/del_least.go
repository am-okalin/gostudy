package stroperation

import (
	"strings"
)

func DelLeast(input string) string {
	length := len(input)
	m := make(map[rune]int, length)
	for _, r := range input {
		m[r]++
	}

	tmp := 0
	for _, num := range m {
		if tmp == 0 || tmp > num {
			tmp = num
		}
	}

	for r, num := range m {
		if num == tmp {
			delete(m, r)
		}
	}

	var builder strings.Builder
	builder.Grow(length)
	for _, r := range input {
		if _, ok := m[r]; ok {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
