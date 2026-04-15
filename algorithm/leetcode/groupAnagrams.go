package leetcode

import "sort"

func groupAnagrams(strs []string) [][]string {
	// 空判断
	if len(strs) == 0 {
		return [][]string{}
	}

	// 初始化
	var arrs = make([][]string, 0)
	var maps = make(map[string][]string)

	// 循环
	for _, str := range strs {
		// 排序找出 zmywc.key
		rs := []rune(str)
		sort.Slice(rs, func(i, j int) bool {
			return rs[i] < rs[j]
		})

		// 加入maps
		rss := string(rs)
		if maps[rss] == nil {
			maps[rss] = make([]string, 0)
		}
		maps[rss] = append(maps[rss], str)
	}

	// maps 转 arrs
	for _, m := range maps {
		arrs = append(arrs, m)
	}

	return arrs
}
