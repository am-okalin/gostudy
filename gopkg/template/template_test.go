package template

import (
	"bytes"
	"testing"
	"text/template"
)

type Inventory struct {
	ShopName string `template:"STORE_NAME"`
	Count    uint
}

func Test1(t *testing.T) {
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Delims("*|", "|*").
		Parse("*|.Count|* items are made of *|.STORE_NAME|*")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, sweaters)
	if err != nil {
		panic(err)
	}
	t.Log(buf.String())
}
