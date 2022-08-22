package graph

var res [][]string

//eQueen queens皇后的放置方法, n表示棋盘边长, row表示当前行
func eQueen(queens []int, n, row int) {
	if row == n { // 符合条件的递归出口
		addAnswer(queens, n)
		return
	}
	for col := 0; col < n; col++ {
		if !isOk(queens, row, col) {
			continue
		}
		queens[row] = col
		eQueen(queens, n, row+1)
		queens[row] = -1
	}
}

//isOk 是否达到枝剪条件
func isOk(queens []int, row, col int) bool {
	m := row + col
	d := row - col
	for r := 0; r < row; r++ {
		if queens[r] == col {
			return false
		}
		if r+queens[r] == m {
			return false
		}
		if r-queens[r] == d {
			return false
		}
	}
	return true
}

//addAnswer 将结果集转化为字符串展示
func addAnswer(queens []int, n int) {
	board := []string{}
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = '.'
		}
		row[queens[i]] = 'Q'
		board = append(board, string(row))
	}
	res = append(res, board)
}

//todo::改为循环模式
