package graph

import "testing"

func TestKnapsack(t *testing.T) {
	p := NewPack(9, 5, []int{2, 2, 4, 6, 3})

	p.knapsack(0, 0)

	t.Log(p.maxW)
}
