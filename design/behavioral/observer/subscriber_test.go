package observer

import "testing"

func Test1(t *testing.T) {
	shirtItem := Publisher{}

	observerFirst := &customer{id: "abc@gmail.com"}
	observerSecond := &customer{id: "xyz@gmail.com"}

	shirtItem.Subscribe(observerFirst)
	shirtItem.Subscribe(observerSecond)
	shirtItem.Unsubscribe(observerFirst)

	shirtItem.Notify()
}
