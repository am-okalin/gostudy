package leetcode

import (
	"fmt"
	"log"
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

func ExampleGroupAnagrams() {
	length := groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"})
	log.Fatal(length)
	// Output:
	// {"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}}
}
