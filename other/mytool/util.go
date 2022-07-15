package mytool

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func appendIfAbsent(s []string, t ...string) []string {
	for _, t1 := range t {
		var contains bool
		for _, s1 := range s {
			if s1 == t1 {
				contains = true
				break
			}
		}
		if !contains {
			s = append(s, t1)
		}
	}
	return s
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
