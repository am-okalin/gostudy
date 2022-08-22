[TOC]
## 入门
### 安装
- 使用goland安装: `File`-`Settings`-`Go`-`GOROOT`-`Download...`
- linux安装: 访问[go_download](https://go.dev/dl/)进行下载安装
```shell
wget https://go.dev/dl/go1.18.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz
cat /usr/local/go/VERSION

# 配置环境变量，更改全局goenv
vim /etc/profile.d/go.sh
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
source /etc/profile.d/go.sh

# 开启gomod 包管理器
go env -w GO111MODULE=on

# 设置代理 国内不配置代理会导致`go get/download/install`下载慢或失败
go env -w GOPROXY=https://proxy.golang.com.cn,direct

# 不走 proxy 的私有仓库或组，多个逗号相隔(可选)，若远端是git仓库则使用git的ssh连接
go env -w GOPRIVATE=*.corp.example.com,github.com/org_private
```

#### helloword
- 初始化项目
```shell
mkdir gostudy && cd gostudy
go mod init gostudy
mkdir cmd && cd cmd
vim main.go
# todo...
go build main.go
./main
```
- helloword
```go
package main //当前包名
import "fmt" //引入的包

func main() {
	//使用fmt包的Println方法
	fmt.Println("hello word")
}
```

#### `代码包`与`源码文件`
- 只有`mian`包的`main()`方法才可以作为`可执行程序`的入口，包含`main`的文件叫做`命令源码文件`
- `命令源码文件`要放在`cmd`目录下，若一个项目有多个`命令源码文件`则放在`cmd`下的多个目录
- GO源码以`代码包`为基本组织单位。`代码包`与系统目录一一对应。目录可以有子目录，`代码包`也可以有`子包`。
- 若复杂程序的`程序实体`都放在`命令源码文件`中，会降低可维护性，可读性。因此我们有必要将`逻辑代码`与`初始化代码`进行包级别的隔离。
- e.g. 自定义`hello $name`的方法并在`main()`中进行调用进行调用

```shell
├── cmd
│   └── main.go
├── go.mod
└── hello
    └── hello.go
```
```go
package main

import "gostudy/hello"

func main() {
	//使用hello包的SayHello方法
	hello.SayHello("downing")
}
```
```go
package hello

import "fmt"

func SayHello(name string) {
	//使用fmt包的Printf方法
	fmt.Printf("hello %s\n", name)
}
```

- 类似`hello`包中的所有文件都被称为`库源码文件`，用于集中放置各种待被使用的`程序实体`(全局常量、全局变量、接口、结构体、函数等等)
- 代码包可被代码包引用，若包A依赖包B 则 B 是 A 的`依赖包`，A 是 B 的`触发包`。但禁止代码包间的循环引用

#### 简单测试
- 复杂程序下，我们无法保证`库源码文件`的功能和性能，因此需要引入测试。而存放`测试方法`的文件就是`测试源码文件`。它的文件名通常是`xx_test.go`。

```go
package hello

import "testing"

func TestSayHello(t *testing.T) {
	SayHello("gary")
}
```
```shell
# -run 执行名为xxx的测试函数 -v显示详细测试过程
go test -v -run=^TestSayHello$ gostudy/hello
```

### 类型
- 按底层结构可分为`值类型`与`引用类型`
	+ `值类型`包括所有基本数据类型，数组，结构体
	+ `引用类型`包括指针、`slice`、`map`、`channel`、`func`、`interface`
- 定义`值类型`和`引用类型`变量时，`值类型`会默认分配内存，`引用类型`初始化后才分配内存。


#### 类型定义与声明
- GO静态强类型，编译型语言。在声明变量时必须指定它们的类型
```go
//TestDeclaration 类型声明
func TestDeclaration(t *testing.T) {
	var str1 string
	str2 := ""          //隐式的指定了类型为string
	str3 := new(string) //返回指针
	t.Log(str1, str2, str3, *str3)
}
```

#### 函数
- `函数`是独立的程序实体。我们可以声明有名字的`函数`，也可以声明没名字的`函数`，还可以把它们当做普通的值传来传去。
- 我们能把具有相同签名的`函数`抽象成独立的`函数类型`(如选项模式中的选项)，以作为一组输入、输出（或者说一类逻辑组件）的代表。
- 函数的形参传递的是值的拷贝
- `init()`初始化函数，引入外部包会优先执行`init()`函数对外部包执行一些初始化的操作
- 函数可有多返回值如os包的`File.Write`的签名为`func (file *File) Write(b []byte) (n int, err error)`

#### 数组`array`与切片`slice`
- 数组是定长的(一旦声明长度便不可更改)，切片是可变长的(长度不够时，可进行扩容)
- 数组是值类型；切片是引用类型，它可被看作为指向数组的指针
- 切片除了已知的3种类型声明外，还可用make进行声明
```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

#### 字典类型`map`
- 张三考了20分，李四考了50分，gary考了100分... 老师想看一张总体的分数表
- `map`底层是hash表，可快速获取键值对，时间复杂度为Ｏ(1)，它不存在容量，只有长度，也可用make进行初始化
```go
//Test1 类型声明
func Test1(t *testing.T) {
	ma := map[int]bool{0: true, 1: true} //go中的集合写法
	mb := make(map[string]int, 0)
	t.Log(ma, mb)
}
```

#### 结构体`struct`
- `struct`是一种聚合的数据类型，是由零个或多个任意类型的值聚合成的实体。每个值称为`struct`的成员。
	+ 可类比其他语言的`class`,如人由姓名年龄组成，那么人这个结构体的成员就是姓名和年龄
	+ 结构体还可以表示更复杂的数据结构，如链表，树，图用结构体表示是相当清晰的
- `方法`与`函数`不同，它需要有名字，不能被当作值来看待，最重要的是，它必须隶属于某个自定义类型。`方法`所属的类型会通过其声明中的接收者`receiver`声明体现出来。
- 要注意区分指针方法和值方法，只有指针方法才可以对`receiver`进行修改
- `struct`中的语法糖: 嵌入字段，以及多层嵌入时的屏蔽

#### 接口`interface`
- GO中很重要的一个思想是面向接口编程，如fmt中的string方法
- 接口可以多组合和内嵌，如io包的接口
