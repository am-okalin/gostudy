package reflection

import (
	"reflect"
	"testing"
)

type S struct {
	F string `species:"gopher" color:"blue"`
}

func TestTag(t *testing.T) {
	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	name, flag := st.FieldByName("F")
	t.Log(name.Tag, flag)
	t.Log(field.Tag.Get("color"), field.Tag.Get("species"))
}

func TestValue(t *testing.T) {
	s := S{F: "downing_test"}
	v := reflect.ValueOf(s).FieldByName("F")
	if v.Kind() != reflect.String {
		t.Error("类型不为string")
	}
	t.Log(v.String())
}

func TestName(t *testing.T) {
	rt := reflect.TypeOf(S{})
	n := rt.Field(0).Name
	t.Log(n)

	rv := reflect.ValueOf(S{})
	n = rv.Type().Name()
	t.Log(n)
}
