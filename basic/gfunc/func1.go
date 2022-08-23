package gfunc

import "fmt"

func init() {
	fmt.Println("进行初始化操作...")
}

var f1 = func() {
	fmt.Println("this is f1")
}

//eat 独立的函数类型 可作为输入输出
type eat func(food string) (weight int)

func getEatFruitFun() eat {
	return func(fruit string) int {
		fmt.Printf("eat %s\n", fruit)
		return 1
	}
}

func getEatMeatFun() eat {
	return func(meat string) int {
		fmt.Printf("eat %s\n", meat)
		return 2
	}
}
