package divideandconquer

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	arr := []int{7, 8, 3, 4, 9, 5}
	mergeSort(arr, 0, 5)
	t.Log(arr)
}
