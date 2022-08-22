package divideandconquer

import "testing"

func TestBinarySearchRecursive(t *testing.T) {
	var a = []int{1, 3, 5, 6, 8}
	if bs(a, 0, len(a), 8) != 4 {
		t.Fatal("not found")
	}
}
