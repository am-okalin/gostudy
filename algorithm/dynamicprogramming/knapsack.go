package dynamicprogramming

//knapsack1 动态规划解决01背包
func knapsack1(w, n int, weight []int) int {
	// 创建状态转移表
	var status = make([][]bool, n)
	for i := range status {
		status[i] = make([]bool, w+1)
	}
	status[0][0] = true
	if weight[0] <= w {
		status[0][weight[0]] = true
	}

	for i := 1; i < n; i++ { //动态规划状态转移
		for j := 0; j <= w; j++ {
			if status[i-1][j] == true { //保持上一行状态
				status[i][j] = status[i-1][j]
			}
		}
		for j := 0; j <= w-weight[i]; j++ {
			if status[i-1][j] == true { //找到最后一个重量后放入i
				status[i][j+weight[i]] = true
			}
		}
	}

	// 找到最大重量
	for i := w; i > 0; i-- {
		if status[n-1][i] == true {
			return i
		}
	}
	return 0
}

func knapsack2(w, n int, weight []int) int {
	var status = make([]bool, w+1)
	status[0] = true
	if weight[0] <= w {
		status[weight[0]] = true
	}
	for i := 1; i < n; i++ {
		for j := w - weight[i]; j > 0; j-- {
			if status[j] == true {
				status[j+weight[i]] = true
			}
		}
	}
	for i := w; i > 0; i-- {
		if status[i] == true {
			return i
		}
	}

	return 0
}
