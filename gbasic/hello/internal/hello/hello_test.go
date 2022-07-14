package hello

import "testing"

//以_test结尾的文件就是测试源码文件
//在编译的时候会自动跳过测试源码文件

func Test1(t *testing.T) {
	SayHello("gary")
}

//TestSayHello 测试方法
func TestSayHello(t *testing.T) {
	type args struct {
		name string
	}
	//以切片对象方式存在是为了测试多个用例
	tests := []struct {
		name string
		args args
	}{
		{
			name: "say bob",
			args: args{
				name: "bob",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SayHello(tt.args.name)
		})
	}
}
