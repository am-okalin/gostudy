package reflection

import (
	"errors"
	"reflect"
)

var ErrCannotBeSet = errors.New("this variable cannot be set")

type User struct {
	Name  string `species:"gopher" color:"blue"`
	Hobby string
	Age   int
}

func setValue(to, from reflect.Value) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if s, ok := p.(string); ok {
				err = errors.New(s)
			}
		}
	}()

	//CanSet 判断to是否可执行set函数
	if !to.CanSet() {
		return ErrCannotBeSet
	}

	//Set 设置to变量为from的值,
	to.Set(from)
	return nil
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}
