package bytes

import (
	"bytes"
	"testing"
)

//合并[]byte
func TestMerge(t *testing.T) {
	str := "01"
	byt := []byte(str)

	//1.使用append
	a := append([]byte{}, byt...)
	t.Log(string(a))

	//2.使用bytes.Buffer
	var buf bytes.Buffer
	buf.Write(byt)
	b := buf.Bytes()
	t.Log(string(b))

	//3.使用bytes.Join(s [][]byte, sep []byte) []byte
	s := [][]byte{byt, []byte("23"), []byte("45")}
	c := bytes.Join(s, []byte("-"))
	t.Log(string(c))
}
