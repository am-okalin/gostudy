package io

import (
	"strings"
	"testing"
)

func TestByte(t *testing.T) {
	var a []byte
	b := new([]byte)
	c := make([]byte, 8)

	t.Log(a, b, c, nil)

	if a == nil {
		t.Log("[]byte == nil")
	}
	if b == nil {
		t.Log("&[]byte == nil")
	}
	if c == nil {
		t.Log("c == nil")
	}
}

func TestCopyBuffer(t *testing.T) {
	r1 := strings.NewReader("first reader\n")
	//r2 := strings.NewReader("second reader\n")
	buf := make([]byte, 8)

	nr, err := r1.Read(buf)
	if err != nil {
		return
	}

	t.Log(nr)
}