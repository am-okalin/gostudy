package a

import (
	"basic/internal/cycleimport/b"
	"fmt"
)

func A1(str string) {
	fmt.Println(str, b.Bads)
}
