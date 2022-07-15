package sync

import (
	"testing"
)

func TestConst(t *testing.T) {
	const (
		mutexLocked = 1 << iota // mutex is locked
		mutexWoken
		mutexStarving
		mutexWaiterShift = iota
	)
	t.Log(mutexLocked)
	t.Log(mutexWoken)
	t.Log(mutexStarving)
	t.Log(mutexWaiterShift)
}
