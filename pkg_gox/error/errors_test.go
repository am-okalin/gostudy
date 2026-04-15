package error

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

// 哨兵
var ErrNotFound = errors.New("not found")

// 自定义错误类型
type QueryError struct {
	Query string
	Msg   string
	Err   error
}

func (e QueryError) Error() string {
	return fmt.Sprintf("{msg:%v; err:%v}", e.Msg, e.Err)
}

func (e QueryError) Unwrap() error {
	return e.Err
}

// TestIs 类型判断
func TestIs(t *testing.T) {
	// 使用 直等 判断哨兵错误
	err := ErrNotFound
	if err == ErrNotFound {
		t.Logf("the err is   ErrNotFound")
	}

	// 使用 errors.Is 判断错误链中是否存在ErrNotFound
	err = fmt.Errorf("wrap error: %w", ErrNotFound)
	if errors.Is(err, ErrNotFound) {
		t.Logf("the err wrap ErrNotFound")
	}
}

// TestAs 类型断言
func TestAs(t *testing.T) {
	// error可接受内建类型的值，或实现了error接口的自定义类型
	if ae1, ok := GetQueryError().(QueryError); ok {
		t.Logf("%+v", ae1)
	}

	// error也可接受自定义类型的指针，注意断言时必须正确书写
	if ae2, ok := GetQueryErrorPtr().(*QueryError); ok {
		t.Logf("%+v", ae2)
	}

	// target在错误链中存在? 若存在则将`其`值赋值给target(即便有值)
	err := fmt.Errorf("wrap error: %w", GetQueryError())
	target := QueryError{Err: io.EOF, Msg: "输入结束"}
	if errors.As(err, &target) {
		t.Logf("%T, %+v", err, err)
		t.Logf("%T, %+v", target, target)
	}
}

func GetQueryError() error {
	return QueryError{Err: ErrNotFound, Msg: "无指针错误"}
}

func GetQueryErrorPtr() error {
	return &QueryError{Err: ErrNotFound, Msg: "有指针错误"}
}
