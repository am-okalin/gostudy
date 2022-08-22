package goruntine

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

//TestStack panic时会打印出goroutine的ID, 状态，函数，调用栈的信息。
func TestStack(t *testing.T) {
	panic("get stack info")
}

//runtime.Stack(buf, all),buf是保存信息的字符串数字
//all表示是否要打印全部的堆栈信息。all=true就会先调用stopTheWorld
func TestGetStackInfo(t *testing.T) {
	var (
		size = 64 * 1024
		all  = false
	)
	buf := make([]byte, size)
	runtime.Stack(buf, all)
	t.Logf("%s", string(buf))
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func TestGoID(t *testing.T) {
	id := GoID()
	t.Log(id)
}
