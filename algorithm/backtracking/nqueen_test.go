package graph

import "testing"

func TestInitBoard(t *testing.T) {
	n := 8
	var queens = make([]int, n)
	for i := range queens {
		queens[i] = -1
	}

	eQueen(queens, n, 0)

	for _, re := range res {
		t.Log(re)
	}
}
