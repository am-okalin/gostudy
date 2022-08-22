package recursive

//Fib1 递归
func Fib1(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 1
	}
	return Fib1(n-1) + Fib1(n-2)
}

//Fib2 尾递归
func Fib2(n, first, second int) int {
	if n == 1 {
		return second
	}

	return Fib2(n-1, second, first+second)
}

//Fib3 尾递归 抽象版-记录每个迭代结果集
func Fib3(n int, sum *[]int) int {
	l := len(*sum)

	if n < 2 {
		return (*sum)[l-1]
	}

	// do something
	*sum = append(*sum, (*sum)[l-1]+(*sum)[l-2])

	return Fib3(n-1, sum)
}

//Fib4 递归+备忘录
func Fib4(n int, memo map[int]int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 1
	}

	if value, ok := memo[n]; ok {
		return value
	}

	memo[n] = Fib1(n-1) + Fib1(n-2)

	return memo[n]
}

//Fib5 循环
func Fib5(n int) int {
	sum := []int{0, 1}

	for i := 2; i <= n; i++ {
		sum = append(sum, sum[i-1]+sum[i-2])
	}

	return sum[n]
}
