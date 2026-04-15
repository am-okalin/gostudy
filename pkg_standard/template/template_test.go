package template

import (
	"bytes"
	"testing"
	"text/template"
)

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

func TestRender(t *testing.T) {
	str := "*|SHOP_DOMAIN|*/discount/*|INVITEE_CODE|*"
	str = Render(str, &Inventory{
		ShopName: "wtn_shop_domain",
		Count:    10,
	}, "*|", "|*")
	t.Log(str)
}
