package stack

import (
	"fmt"
	"testing"
)

//matchBrackets 匹配括号
func matchBrackets(brackets string) bool {
	m := map[rune]rune{')': '(', ']': '[', '}': '{'}
	rs := stack{cap: 32, data: make([]any, 0, 32)}

	for _, b := range brackets {
		//栈报错返回false
		if rs.Err() != nil {
			return false
		}

		//非括号不处理
		_, ok := m[b]

		//左括号则入栈
		if !ok && b != ' ' {
			rs.Push(b)
			fmt.Printf("入栈:%s\n", rs.String())
		}

		//右括号则出栈
		if ok {
			val := rs.Pop()
			//出栈字符 不是 右括号的匹配括号
			if val != m[b] {
				return false
			}
			fmt.Printf("出栈:%s->%c%c\n", rs.String(), val, b)
		}

	}
	return rs.top == 0
}

func TestMatchBrackets(t *testing.T) {
	b1 := matchBrackets("{() }[]")
	t.Log(b1)
	b2 := matchBrackets("(]")
	t.Log(b2)
}
