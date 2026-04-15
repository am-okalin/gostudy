package regexp

import (
	"regexp"
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	url := "http://127.0.0.1:7150/api/buildings/ao8bmbh24dfi/maps/{mid}/fence/{fid}"
	expr := "{(.+?)\\}"

	// 解析正则表达式
	re, err := regexp.Compile(expr)
	if err != nil {
		t.Error(err)
		return
	}

	// 第二个参数n小于0返回全部匹配，大于0则返回前n个
	list := re.FindAllString(url, -1)
	t.Log(list)

	rep := re.ReplaceAllStringFunc(url, strings.ToUpper)
	t.Log(rep)

	//把匹配的所有字符a替换成b
	rep2 := re.ReplaceAllString(url, "b")
	t.Log(rep2)
}

func Test2(t *testing.T) {
	a := strings.ReplaceAll("assets/cmp/referral/popcat.png", "/", "_")
	t.Log(a)
}
