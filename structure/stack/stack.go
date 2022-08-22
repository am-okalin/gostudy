package stack

import (
	"errors"
	"strings"
)

type StackInterface interface {
	Len() int
	Err() error
	Pop() any
	Push(v any)
	String() string
}

type stack struct {
	top  int   //头指针 0表示栈无数据
	cap  int   //栈容量 等同于 len(data)
	err  error //错误信息
	data []any //数据域
}

func (s *stack) Len() int {
	return s.cap
}

func (s *stack) Err() error {
	return s.err
}

func (s *stack) Pop() any {
	//有错误直接弹出
	if s.err != nil {
		return nil
	}
	//空栈不可取
	if s.top <= 0 {
		s.err = errors.New("stack is empty")
		return nil
	}
	//弹出数据
	s.top--
	val := s.data[s.top]
	s.data = s.data[:s.top]
	return val
}

func (s *stack) Push(v any) {
	//有错误直接弹出
	if s.err != nil {
		return
	}
	//满栈不可入
	if s.top >= s.cap-1 {
		s.err = errors.New("stack is empty")
		return
	}
	//压入数据
	s.top++
	s.data = append(s.data, v)
	return
}

func (s *stack) String() string {
	b := strings.Builder{}
	for i := 0; i < s.top; i++ {
		b.WriteRune(s.data[i].(rune))
	}
	return b.String()
}
