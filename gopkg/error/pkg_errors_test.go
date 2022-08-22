package error

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"testing"
)

func TestPkgErr(t *testing.T) {
	//使用errors.New 或 errors.Errorf
	e0 := errors.Errorf("%v", io.EOF)
	//包装错误但不记录堆栈
	e1 := errors.WithMessagef(e0, "e1")
	//包装错误及堆栈
	e2 := errors.Wrap(e1, "e2")
	//获取根因
	e3 := errors.Cause(e2)
	//打印错误
	fmt.Printf("%T %v\n", e1, e1)
	fmt.Printf("%T %v\n", e2, e2)
	fmt.Printf("%T %v\n", e3, e3)
	fmt.Printf("%+v\n", e2)
}
