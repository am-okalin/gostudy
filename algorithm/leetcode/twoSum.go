package leetcode

//twoSum 两数之和 从nums中找出和为target的两个数
func twoSum(nums []int, target int) []int {
	m := make(map[int]int, len(nums))
	for i, num := range nums {
		diff := target - num
		if _, ok := m[diff]; ok {
			return []int{i, m[diff]}
		}
		m[num] = i
	}
	return nil
}
