package error

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

//哨兵
var ErrNotFound = errors.New("not found")

type QueryError struct {
	Query string
	Msg   string
	Err   error
}

func (e QueryError) Error() string {
	return fmt.Sprintf("query_error:%v", e.Msg)
}

func (e QueryError) Unwrap() error {
	return e.Err
}

func Test1(t *testing.T) {
	//false 两变量的地址不一样
	e1 := errors.New("not found")
	fmt.Println(ErrNotFound == e1)
	//error可接受内建类型的值，或实现了error接口的自定义类型
	assertQE(QueryError{Err: ErrNotFound})
	//error也可接受自定义类型的指针，注意断言时必须正确书写
	assertQEP(&QueryError{Err: e1})
}

func assertQE(err error) {
	e, ok := err.(QueryError)
	fmt.Println(e, ok)
}

func assertQEP(err error) {
	e, ok := err.(*QueryError)
	fmt.Println(e, ok)
}

func Test2(t *testing.T) {
	//传参不能是指针类型，否则在反射包中就无法成功
	assertError(QueryError{Err: ErrNotFound})
}

func assertError(err error) {
	//错误判断,只要错误链中存在ErrNotFound就返回true
	if errors.Is(err, ErrNotFound) {
		fmt.Println("the err wrap ErrNotFound")
	}

	//错误断言,As的第二个参数target必须是指向错误的指针
	target := &QueryError{Err: io.EOF}
	if errors.As(err, target) {
		//此时target已经指向了err的内容
		fmt.Println(target.Err)
	}
}

func Test3(t *testing.T) {
	//fmt.Errorf用fmt.wrapError包装了err
	err := fmt.Errorf("wrap error: %w", ErrNotFound)
	if errors.Is(err, ErrNotFound) {
		fmt.Printf("%T, %v", err, ErrNotFound)
	}
}
