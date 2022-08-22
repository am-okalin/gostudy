package queue

type Node int

type QueueInterface interface {
	IsEmpty() bool
	IsFull() bool
	Length() int
	Capacity() int
	Front() Node
	Rear() Node
	EnQueue(value Node) bool
	DeQueue() bool
}
