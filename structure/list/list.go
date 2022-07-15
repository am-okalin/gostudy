package list

import "strings"

type Elem struct {
	next  *Elem
	Value any
}

func (e *Elem) String() string {
	return e.Value.(string)
}

type List struct {
	head *Elem
}

//get 获取第i个位置的元素
func (l *List) get(i int) *Elem {
	//取第i个元素，只用快指针
	slow, fast := l.head, l.head
	if i >= 0 {
		for j := 0; j < i && fast != nil; j++ {
			fast = fast.next
		}
		return fast
	}
	//取倒数第i个元素，要用快慢指针
	for j := i; j < 0; j++ {
		if fast = fast.next; fast == nil {
			return nil
		}
	}
	//fast为nil时，slow就是倒数第i个元素
	for fast != nil {
		fast = fast.next
		slow = slow.next
	}
	return slow
}

//insert 在第i个位置插入e
func (l *List) insert(e, at *Elem) {
	e.next = at.next
	at.next = e
}

//remove 删除at后的元素
func (l *List) removeAt(at *Elem) bool {
	if at.next == nil {
		return false
	}

	e := at.next
	at.next = e.next
	e.next = nil
	return true
}

//reverse 链表反转
func (l *List) reverse() {
	if l.head == nil || l.head.next == nil {
		return
	}
	left := l.head
	curr := l.head.next
	for curr != nil {
		right := curr.next
		curr.next = left
		left = curr
		curr = right
	}
	l.head.next = nil
	l.head = left
}

//hasCycle 检测是否有环
func (l *List) hasCycle() bool {
	if l.head != nil {
		slow, fast := l.head, l.head
		for fast != nil && fast.next != nil {
			slow = slow.next
			fast = fast.next.next
			if fast == slow {
				return true
			}
		}
	}
	return false
}

//middle 返回中间元素
func (l *List) middle() *Elem {
	slow, fast := l.head, l.head
	if l.head != nil {
		for fast != nil && fast.next != nil {
			slow = slow.next
			fast = fast.next.next
		}
	}
	return slow
}

func (l *List) String() string {
	s := strings.Builder{}
	e := l.head
	s.WriteString(e.Value.(string))
	for e.next != nil {
		e = e.next
		s.WriteString(" -> ")
		s.WriteString(e.Value.(string))
	}
	return s.String()
}
