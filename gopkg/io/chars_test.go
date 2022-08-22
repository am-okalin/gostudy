package io

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	str1 := "壹，二,3\\four"
	fmt.Println(str1)

	//rune int32 4x8-1位
	runes := []rune(str1)
	fmt.Println(runes[0], string(runes[0]))
	fmt.Printf("%v%c\n", runes[1], runes[1])

	//byte uint8 1x8位
	bytes := []byte(str1)
	fmt.Printf("%v\n", bytes)
	fmt.Printf("%s\n", bytes[0:6])
}
