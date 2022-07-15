package dynamicprogramming

import "testing"

func TestKnapsack1(t *testing.T) {
	maxW := knapsack1(9, 5, []int{2, 2, 4, 6, 3})
	t.Log(maxW)
}

func TestKnapsack2(t *testing.T) {
	maxW := knapsack2(9, 5, []int{2, 2, 4, 6, 3})
	t.Log(maxW)
}
