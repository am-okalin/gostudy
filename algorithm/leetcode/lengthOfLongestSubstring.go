package leetcode

func lengthOfLongestSubstring(s string) int {
	charMap := make(map[byte]int) // 记录字符在窗口中出现的次数
	left := 0
	maxLen := 0

	for right := 0; right < len(s); right++ {
		char := s[right]
		charMap[char]++

		// 如果当前字符重复，移动左指针缩小窗口
		for charMap[char] > 1 {
			charMap[s[left]]--
			left++
		}

		// 更新最大长度
		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
		}
	}
	return maxLen
}
