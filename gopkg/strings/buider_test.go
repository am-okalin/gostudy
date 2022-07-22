package strings

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestStr(t *testing.T) {
	var builder1 strings.Builder
	builder1.WriteString("wtn_test")
	builder1.WriteString("hahahah")
	s := builder1.String()
	t.Log(s)
}

func TestBuilder(t *testing.T) {
	f1 := func(b1 strings.Builder) {
		time.Sleep(1 * time.Second)
		fmt.Println("f1_builder: ", b1.String())
	}
	f2 := func(b2 *strings.Builder) {
		time.Sleep(1 * time.Second)
		fmt.Println("f2_builder: ", b2.String())
	}

	var builder strings.Builder
	go f1(builder)
	go f2(&builder)
	//写入单个字符
	builder.WriteRune(' ')
	builder.WriteByte(' ')
	builder.Write([]byte("some_str"))
	builder.WriteString("some_str")
	fmt.Println("main_builder: ", builder.String())
	f2(&builder)
	time.Sleep(3 * time.Second)
}
