package queue

type QueueArr struct {
	data       *[]Node //用切片存储队列数据，出队时无需清除元素内容
	capacity   int     //切片所指数组的容量，也指队列总长度
	head, tail int     //逻辑指针，物理指针为 切片所指数组的起始地址+Node位宽×head
}

func NewQueueArr(capacity int) *QueueArr {
	data := make([]Node, capacity, capacity)
	// 逻辑指针指向-1表示空
	return &QueueArr{data: &data, capacity: capacity, head: -1, tail: -1}
}

func (q *QueueArr) Capacity() int {
	return q.capacity
}

func (q *QueueArr) IsEmpty() bool {
	return q.tail == -1 && q.head == -1
}

func (q *QueueArr) IsFull() bool {
	return (q.tail+1)%q.capacity == q.head
}

func (q *QueueArr) Length() int {
	if q.IsEmpty() {
		return 0
	}
	return (q.tail - q.head + q.capacity + 1) % q.capacity
}

func (q *QueueArr) Front() Node {
	if q.IsEmpty() {
		return -1
	}
	return (*q.data)[q.head]
}

func (q *QueueArr) Rear() Node {
	if q.IsEmpty() {
		return -1
	}
	return (*q.data)[q.tail]
}

func (q *QueueArr) EnQueue(value Node) bool {
	if q.IsFull() {
		return false
	}
	if q.IsEmpty() {
		q.head = 0
	}
	q.tail = (q.tail + 1) % q.capacity
	(*q.data)[q.tail] = value
	return true
}

func (q *QueueArr) DeQueue() bool {
	if q.IsEmpty() {
		return false
	}
	if q.head == q.tail {
		q.head = -1
		q.tail = -1
	} else {
		q.head = (q.head + 1) % q.capacity
	}
	return true
}
