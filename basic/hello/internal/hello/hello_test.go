package hello

import "testing"

//以_test结尾的文件就是测试源码文件
//在编译的时候会自动跳过测试源码文件

func Test1(t *testing.T) {
	SayHello("gary")
}

func Test_sayHello(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test_say_downing", args{name: "downing"}},
		{"test_say_gary", args{name: "gary"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sayHello(tt.args.name)
		})
	}
}
