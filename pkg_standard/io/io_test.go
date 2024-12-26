package io

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	//copy 除了同类型切片复制外还可将string复制到[]byte类型
	src := "零一二三四五六七八九"
	dst := make([]byte, 0, 512)
	dst = append(dst, []byte("abcdefg")...)

	//copy 会将 dst 中的数据覆盖
	//copy 长度等同于len(dst)
	n := copy(dst, src)
	t.Log(n, string(dst))
}

func TestCopyBuffer(t *testing.T) {
	r1 := strings.NewReader("first reader\n")
	//r2 := strings.NewReader("second reader\n")
	buf := make([]byte, 8)

	nr, err := r1.Read(buf)
	if err != nil {
		return
	}

	t.Log(nr)
}

//对 io.ReadAll(r) 方法的测试
func TestReader(t *testing.T) {
	b := make([]byte, 0, 512)
	b = append(b, []byte("ab")...)
	r := strings.NewReader("0123456789")

	// 临时变量b1是为了将r中数据读入b的底层数组
	b1 := b[len(b):cap(b)]
	n, err := r.Read(b1)
	// 底层数组更新后更新b的长度len
	b = b[:len(b)+n]
	if err != nil {
		t.Error(err)
	}
	t.Log(b)
}

func TestSectionReader(t *testing.T) {
	r := strings.NewReader("0123456789")

	sr := io.NewSectionReader(r, 4, 10)

	buf, err := io.ReadAll(sr)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf))
}

func TestLimitedReader(t *testing.T) {
	r := strings.NewReader("0123456789")

	lr := io.LimitReader(r, 6)

	p := make([]byte, 0, 4)

	n, err := lr.Read(p)
	if err != nil {
		return
	}
	t.Log(n)
	t.Log(string(p))

	n, err = lr.Read(p)
	if err != nil {
		return
	}
	t.Log(n)
	t.Log(string(p))
}

func TestPipe(t *testing.T) {
	r, w := io.Pipe()

	go func() {
		// 写操作阻塞
		fmt.Fprint(w, "some io.Reader stream to be read\n")
		w.Close()
	}()

	// 读操作阻塞
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	// Output:
	// some io.Reader stream to be read
}
