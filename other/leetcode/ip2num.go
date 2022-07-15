package leetcode

import "fmt"

func ip2num(ip string) int {
	slice := make([]int, 4, 4)
	_, _ = fmt.Sscanf(ip, "%d.%d.%d.%d", &slice[0], &slice[1], &slice[2], &slice[3])

	num := slice[0]<<24 + slice[1]<<16 + slice[2]<<8 + slice[3]
	return num
}

func num2ip(num int) string {
	s := fmt.Sprintf("%b", num)
	return s
}
