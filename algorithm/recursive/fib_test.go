package recursive

import "testing"

func TestFib1(t *testing.T) {
	var sum []int

	sum = append(sum, Fib1(1))
	sum = append(sum, Fib1(2))
	sum = append(sum, Fib1(3))
	sum = append(sum, Fib1(4))
	sum = append(sum, Fib1(5))

	println(sum)
}

func TestFib2(t *testing.T) {
	sum := Fib2(5, 1, 1)

	println(sum)
}

func TestFib3(t *testing.T) {
	// 尾递归的出口值有参数写入
	sum := Fib3(5, &[]int{1, 1})

	println(sum)
}

func TestFib4(t *testing.T) {
	sum := Fib4(5, map[int]int{})

	println(sum)
}

func TestFib5(t *testing.T) {
	sum := Fib5(5)

	println(sum)
}
