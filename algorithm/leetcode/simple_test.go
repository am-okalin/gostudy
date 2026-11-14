package leetcode

import (
	"fmt"
	"testing"
)

func TestTwoSum(t *testing.T) {
	nums := []int{2, 6, 3, 7}
	target := 9
	sum := twoSum(nums, target)
	fmt.Println(sum)
}

func TestIp2Num(t *testing.T) {
	num2 := ip2num("10.0.3.193")
	ip2 := num2ip(167969729)
	t.Log(num2, ip2)
}

func TestGroupAnagrams(t *testing.T) {
	groups := groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"})
	if len(groups) != 3 {
		t.Fatalf("groupAnagrams() = %v, want 3 groups", groups)
	}
	total := 0
	for _, group := range groups {
		total += len(group)
	}
	if total != 6 {
		t.Fatalf("groupAnagrams() = %v, want 6 elements in total", groups)
	}
}
