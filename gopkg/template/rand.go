package template

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const TagName = "template"

type Inventory struct {
	ShopName string `template:"STORE_NAME"`
	Count    uint
}

func getKey(t reflect.StructField, tagName string) string {
	tag := t.Tag.Get(tagName)
	if tag == "" {
		return t.Name
	}
	return tag
}

// Render 从结构体渲染到变量
func Render(str string, variables interface{}, left, right string) string {
	t, v := reflect.TypeOf(variables), reflect.ValueOf(variables)
	if v.Elem().Kind() == reflect.Struct {
		t, v = t.Elem(), v.Elem()
	}

	if v.Kind() != reflect.Struct {
		err := fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
		log.Fatal(err)
		return str
	}

	for i := 0; i < t.NumField(); i++ {
		key := fmt.Sprintf("%s%s%s", left, getKey(t.Field(i), TagName), right)
		val := fmt.Sprintf("%v", v.Field(i).Interface())
		str = strings.Replace(str, key, val, -1)
	}

	return str
}
