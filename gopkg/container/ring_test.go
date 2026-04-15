package container

import (
	"container/ring"
	"testing"
)

func TestRing(t *testing.T) {
	// 1. ring.New(n) 会创建一个有 n 个未初始化节点的环
	r := ring.New(5)

	// 2. 初始化环中的值
	// 因为 ring 是循环的，我们只需要遍历一次 (r.Len()) 即可填满
	initValue := 100
	r.Do(func(x any) {
		r.Value = initValue
		initValue += 100
	})

	// 验证初始化结果: [100, 200, 300, 400, 500]
	t.Run("Initialization", func(t *testing.T) {
		expected := []int{100, 200, 300, 400, 500}
		i := 0
		r.Do(func(x any) {
			val := x.(int)
			if val != expected[i] {
				t.Errorf("Expected %d, got %d", expected[i], val)
			}
			i++
		})
	})

	// 3. 测试 Move (移动指针)
	// 当前 r 指向第一个元素 (100)。Move(2) 向后移动 2 步，指向第三个元素 (300)
	r = r.Move(2)
	if r.Value.(int) != 300 {
		t.Errorf("Expected cursor to be at 300, got %d", r.Value.(int))
	}

	// 4. 测试 Link (连接新环/插入)
	// 创建一个包含 2 个元素的新环
	r2 := ring.New(2)
	r2.Value = 999
	r2.Next().Value = 888

	// 将 r2 链接到 r 之后
	// 逻辑：r -> r2 -> (r原本后面的部分)
	// 注意：Link 会返回 r2 的最后一个元素，但通常我们关注结构变化
	r.Link(r2)

	// 验证长度变化：原5 + 新2 = 7
	if r.Len() != 7 {
		t.Errorf("Expected length 7, got %d", r.Len())
	}

	// 打印最终结构验证
	var result []int
	r.Do(func(x any) {
		result = append(result, x.(int))
	})
	// 预期顺序：300 (当前r位置) -> 400 -> 500 -> 100 -> 200 -> 999 -> 888
	// 注意：因为我们在 300 处插入，999 和 888 会插在 300 后面
	t.Log("Final Ring Structure:", result)
}
