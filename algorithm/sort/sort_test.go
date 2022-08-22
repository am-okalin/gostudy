package sort

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	arr := []int{13, 4, 8, 2, 3, 45, 5, 2, 1, 3}

	func(arr []int) {
		qsort(arr, 0, len(arr)-1)
	}(arr)

	fmt.Println(arr)
}
