package sort

func qsort(arr []int, start, end int) {
	if start >= end {
		return
	}

	pivot := arr[end]
	slow := start
	for fast := start; fast < end; fast++ {
		if arr[fast] < pivot {
			if slow != fast {
				arr[fast], arr[slow] = arr[slow], arr[fast]
			}
			slow++
		}
	}

	arr[end], arr[slow] = arr[slow], arr[end]

	qsort(arr, start, slow-1)
	qsort(arr, slow+1, end)
}
