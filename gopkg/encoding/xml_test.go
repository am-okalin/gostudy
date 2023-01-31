package encoding

import (
	"encoding/xml"
	"os"
	"testing"
)

const (
	tryout ="./tryout_heyup.xlsx"
)

func Test1(t *testing.T) {
	//获取文件句柄
	file, err := os.OpenFile(tryout, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.
}
