package net

import (
	"syscall"
	"testing"
)

func TestSocket(t *testing.T) {
	var domain int
	var stype int
	var proto int
	socket, err := syscall.Socket(domain, stype, proto)
	if err != nil {
		return
	}
	t.Log(socket)
}

func TestRune(t *testing.T) {

}
