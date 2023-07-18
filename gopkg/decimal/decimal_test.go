package decimal

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

const (
	Pi = "3.1415926"
)

func Test1(t *testing.T) {
	pi, _ := decimal.NewFromString(Pi)

	a := pi.RoundFloor(2).String()

	// 输出结果
	fmt.Println(a)
}
