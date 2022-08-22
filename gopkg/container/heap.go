package container

//IntHeap heap提供了接口，需要自己实现如下方法
type IntHeap []int

//Less 构造的是小顶堆，大顶堆只需要改一下下面的符号
func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Pop() interface{} {
	x := (*h)[h.Len() - 1]
	*h = (*h)[: h.Len() - 1]
	return x
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}