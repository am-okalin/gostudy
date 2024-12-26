package text

import (
	"golang.org/x/text/language"
	"testing"
)

func Test1(t *testing.T) {
	tag := language.Make("el")
	t.Log(tag)
	parse, err := language.Parse("en-UK")
	t.Log(parse, err)

	ja, _ := language.ParseBase("ja")
	jp, _ := language.ParseRegion("JP")
	jpLngTag, _ := language.Compose(ja, jp)
	t.Log(jpLngTag) // prints

	// 制定了无效的标签
	t.Log(language.Compose(language.ParseRegion("AL"))) // prints Und-AL// ja-JP
}
