package reflection

import (
	"fmt"
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
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
}
