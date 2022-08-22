package sync

import (
	"fmt"
	"testing"
)

func TestBehavior(t *testing.T) {
	var list []string
	for i := 0; i < 100; i++ {
		list = append(list, fmt.Sprintf("string %v", i))
	}

	err := Consume(list, 10)
	if err != nil {
		return
	}
}

func TestConsume1(t *testing.T) {
	var list []T
	for i := 0; i < 100; i++ {
		list = append(list, T{
			a: i,
			b: "",
		})
	}

	err := Consume1(list, 10)
	if err != nil {
		return
	}
}
