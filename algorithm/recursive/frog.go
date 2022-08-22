package recursive

// 递归-备忘录递归-尾递归-递推
// 青蛙一次可以跳上1级台阶，也可以跳上2级台阶。求该青蛙跳上一个 n 级的台阶总共有多少种跳法。
func frog1(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	return frog1(n-1) + frog1(n-2)
}

// 带备忘录的自顶向下
func frog2(n int, m map[int]int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}

	if _, ok := m[n]; ok {
		return m[n]
	}
	m[n] = frog2(n-1, m) + frog2(n-2, m)

	return m[n]
}

// 自底向上
func frog3(n int) int {
	m := make(map[int]int)
	for i := 0; i < n; i++ {
		if i == 1 {
			m[i] = 1
		} else if i == 2 {
			m[i] = 2
		} else {
			m[i] = m[i-1] + m[i-2]
		}
	}

	var sum int
	for _, v := range m {
		sum += v
	}

	return sum
}
