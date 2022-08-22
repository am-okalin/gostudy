## 命令行
- `go help <command>`查看`<command>`基础命令的详细帮助文档
- `go help <topic>`查看`<topic>`扩展命令的详细帮助文档
- `govendor`目前已弃用，建议`go mod`。

### environment
- `CGO_ENABLED` 指明cgo工具是否可用的标识
- `GOHOSTARCH` 程序运行环境的编译体系结构(用于交叉编译，默认与 GOARCH 相同)
- `GOHOSTOS` 程序运行环境的操作系统(用于交叉编译，默认与 GOOS 相同)
- `GOARCH` 程序构建环境的目标计算架构
- `GOOS` 程序构建环境的目标操作系统
- `GOMAXPROCS` 用于设置应用程序可使用的处理器个数与核数，即`processor`的数量
- `GOBIN` 表示编译器和链接器的安装位置，默认是 $GOROOT/bin
- `GOEXE` 可执行文件的后缀。
- `GOENV` go环境配置文件的位置。
- `GOROOT` go编译器源码目录
- `GOPATH` go工作空间 默认为`$HOME/go` 该路径下有三个子目录
    + `src` 存放源码，如`.go` `.c` `.h` `.s`
    + `pkg` 编译时生成的中间文件(包对象)如`.a`,`install`会将编译好的包直接从这里引用,不必再次构建
    + `bin` 编译后生成的可执行文件，`PATH=$PATH:/usr/local/go/bin`
- `GOTOOLDIR` Go工具目录的绝对路径。
- `GOPROXY` 例`https://goproxy.cn,direct,EOF`，表示代理地址列表，当前代理地址40X时自动访问下一个，`direct`表示源地址抓取(如`github`),遇到`EOF`时终止并抛错`invalid version: unknown revision...`
- `GO111MODULE` 
    + `on` 使用Go Modules,go 会忽略 $GOPATH和 vendor文件夹,只根据go.mod下载依赖。
    + `off` 会查找 vendor目录和 $GOPATH来查找依赖关系，也就是继续使用`GOPATH模式`
    + `auto` (默认模式)根据当前项目目录下是否存在 go.mod文件或 $GOPATH/src之外并且其本身包含go.mod文件时才会使用新特性 Go Modules模式

### go env
- `go env [-json] [-u] [-w] [VAR ...]` 设置/查看环境变量，指定了`[VAR ...]`则只打印/设置选中的环境变量的值
- `-json` 以`json`格式打印环境变量
- `-u` 重置变量`VAR [...]`为初始值，以空格分隔。撤销`GOENV`中对应变量的修改
- `-w` 修改变量`VAR=VALUE [...]`，以空格分隔。修改内容会写入`GOENV`中，重启后也会生效
- 系统环境变量的值会覆盖`go env`的配置

### build/install/run/test flag公共参数
- `-n` 打印编译时会用到的所有命令，但不运行
- `-x` 打印编译时会用到的所有命令
- `-a` 强制对所有`包`重新构建
- `-v` 编译时显示包名
- `-p n` 开启并发编译，默认为cpu的逻辑核数
- `-race` 开启竞态条件检测
- `-work` 打印出编译时生成的`临时文件`的路径，并在编译结束时保留它。
- `-gcflags` 给go编译器传入参数, 使用`go tool compile`的参数
- `-ldflags` 给go链接器传入参数, 使用`go tool link `的参数

### go build
- 编译`命令源码文件`会在当前目录下生成`可执行文件`，不可编译多个`命令源码文件`，因为多个`main`函数在编译时会抛出重复定义错误 
- 编译`库源码文件`/`测试源码文件`时，日做检查性编译，不输出任何结果文件。
- 编译A包时会先检测A的`依赖包`，若其编译状态不是最新，则先编译`依赖包`然后编译A
- 不加参数代表编译当前目录对应代码包
- `-o` 指定输出文件名称
- `-i` 安装`依赖包`。即产生`依赖包`的归档文件至`$GOPATH/pkg/...`下

### go install
- `install`只比`build`多做了一件事: 安装编译后的`结果文件`到`指定目录`
- `命令源码文件`的编译生成的`可执行文件`保存在`$GOROOT/bin`下
- `库源码文件`的编译生成的`归档文件`保存在`$GOROOT/pkg`下

### go run
- `go run`在`临时目录`中编译并运行`命令源码文件`
- 若命令源码文件有参数如`go run showds.go -p ~/golang/goc2p`等于编译后执行`./showds -p=~/golang/goc2p`

### go get
- `go get` 获取`vcs`的(远程)包至`$GOPATH`，然后执行`go install`构建并安装它
- `-d` 只执行下载，不执行安装(`install`)。
- `-t` 下载并安装指定的代码包中的测试源码文件中依赖的代码包。
- `-u` 更新已有代码包及其依赖包。
- `-f` 与`-u`一起用时有效。该标记会忽略对`已下载代码包的导入路径`的检查
- `-fix` 下载代码包后先执行`修正`动作，而后再进行编译和安装。
- `-insecure` 允许命令程序使用非安全的scheme（如HTTP）去下载指定的代码包。

### go clean
- go clean删除执行其他命令时产生的一些文件或目录
- 删除`go build`生成的可执行文件
- 删除`go test -c`生成的`xxx.test`文件
- 删除编译时生成的`临时文件`
- `-i` 删除当前包在`$GOROOT/bin`或`$GOROOT/pkg`下生成的文件
- `-r` 删除`当前包`及其所有`依赖包`的上述文件
- `-cache` 删除所有`go build`命令的缓存
- `-testcache` 删除当前包所有的测试结果

### go doc
- 打印附于Go语言程序实体上的文档。可把`程序实体的标识符`作为该命令的参数，按`源码包查找方式`找到包并查看其文档
- `go doc`后不加参数表示查看当前包的所有文档，也可打印标准包如`go doc http.Request`
- `-cmd` 输出`main`包中`可导出的程序实体`的文档
- `-u` 输出`不可导出的程序实体`的文档
- `-c` 区分参数中的大小写
- `-src` 查看`程序实体`的源码
- `-ex`查看`程序实体`示例代码
- `-html`查看HTML格式的文档, 如`godoc -http=:6060`

### go test
- `go test basic pkgtool` 测试多个代码包，执行其中`测试源码文件`
- `go test envir_test.go envir.go` 参数为文件名会生成`虚拟代码包`并执行测试
- `-c` 生成用于运行测试的可执行文件`pkg.test`，但不执行它。
    + `pkg`为被测试代码包的导入路径的最后一个元素的名称
- `-i` 安装/重新安装运行测试所需的依赖包，但不编译和运行测试代码
- `-o` 指定用于运行测试的可执行文件的名称。追加该标记不会影响测试代码的运行。
- `go test`还有很多其他标记，应单独出个`go测试`模块的记录

```shell
# 默认执行每个测试文件的第一个测试函数
go test dir/package
# 清理测试缓存(缓存不会影响测试准确性)
go clean -testcache
# -bench 执行任意名称的性能测试函数
# -run 执行名为空的功能测试函数(即不执行功能测试)
go test -bench=. -run=^$ dir/packages
```

### go list
- `go list [-f format]/[-json] [optoin] [packages]` 展示包相关信息
- `[packages]`不填表示当前目录包，`all`表示所有包
- `-f format`以模板格式打印。模板格式默认为`-f '{{.ImportPath}}'`
- `-json` 以`json`格式打印。不设置此选项则默认为模板格式打印
- `-deps` 将`packages`的`依赖包`也打印出来
- `-e` 忽略错误信息打印
- `go list -m [-u] [-retracted] [-versions] [build flags] [modules]`
    + `-m` 使`go list`列出模块而不是包，仅应用与支持`mod`的包目录下
    + `-u` 添加&列出可用的升级信息。使用`go get -u need-upgrade-package`升级后会将新的依赖版本更新到`go.mod`也可以使用`go get -u`升级所有依赖
- [go list的打印字段的含义](https://www.kancloud.cn/cattong/go_command_tutorial/261354)

### go fmt
- `go fmt [path...]`等同于`gofmt -l -w [path...]`
- `gofmt [flags] [path...]` 格式化包或源文件，`flags`如下
    + `-d` 打印差异输出
    + `-e` 打印所有错误
    + `-l` 打印被格式化的文件名
    + `-r 'pattern -> replacement'` 格式化前载入规则
    + `-s` 尝试简化代码,如`s[a:len(s)]`简化为`s[a:]`
    + `-w` 格式化后的源不输出到`stdout`, 格式化文件时如果失败则回滚至源文件

### go tool fix
- `go fix`是`go tool fix`的简单封装
- `go fix`把指定代码包的所有Go语言源码文件中的旧版本代码修正为新版本的代码

### go tool vet
- `go vet`是`go tool vet`的简单封装
- `go vet`用于检查Go语言源码中静态错误的工具，可使用`-n`/`-x`标记
- `vet`属于Go自带的特殊工具，也是比较底层的命令之一。Go语言自带的特殊工具的存放路径是`$GOROOT/pkg/tool/$GOOS_$GOARCH/`，我们暂且称之为Go工具目录。
- `all` 检查全部。如果有其他检查标记被设置，则命令程序会将此值变为false。默认值为true。
- `asmdecl` 对汇编语言的源码文件进行检查。默认值为false。
- `assign` 检查赋值语句。默认值为false。
- `atomic` 检查代码中对代码包sync/atomic的使用是否正确。默认值为false。
- `buildtags` 检查编译标签的有效性。默认值为false。
- `composites` 检查复合结构实例的初始化代码。默认值为false。
- `compositeWhiteList` 是否使用复合结构检查的白名单。仅供测试使用。默认值为true。
- `methods` 检查那些拥有标准命名的方法的签名。默认值为false。
- `printf` 检查代码中对打印函数的使用是否正确。默认值为false。
- `printfuncs` 需要检查的代码中使用的打印函数的名称的列表，多个函数名称之间用英文半角逗号分隔。默认值为空字符串。
- `rangeloops` 检查代码中对在```range```语句块中迭代赋值的变量的使用是否正确。默认值为false。
- `structtags` 检查结构体类型的字段的标签的格式是否标准。默认值为false。
- `unreachable` 查找并报告不可到达的代码。默认值为false

### Go Modules
- `GO111MODULE` 可配置GoModules开关, `off`关闭，`on`开启，`auto`表示在`$GOPATH/src`下，且没有包含`go.mod`时则关闭`Go Modules`，其他情况下都开启`Go Modules`
- 使用模块时`GOPATH`仅用于存储源码`GOPATH/pkg/mod/`与命令`GOPATH/bin/`
- `go module`根据`go.mod`安装`package`的原則是先拉最新的`release tag`，若无`tag`则拉最新的`commit`。 `golang`会自动生成一个`go.sum`文件来记录`dependency tree`

### `go mod <command> [arguments]` 仅在开启后可用
- `go mod <command> [arguments]` `<command>`如下
    + `download` 下载go.mod 文件中的依赖包
    + `edit` 编辑`go.mod`
    + `graph` 打印模块依赖图
    + `init` 当前目录内初始化`mod`(创建`go.mod`)
    + `tidy` 拉去缺少的模块，移除不用的模块
    + `vendor` 将依赖复制到vendor下(是不是要兼容旧版本的go？)
    + `verify` 验证依赖是否正确
    + `why` 解释为什么需要依赖
- `go.mod`文件语法中提供了以下四个命令
    + `module` 指定包的名字(路径)
    + `exclude` 忽略依赖项模块
    + `require` 指定的依赖项模块 `<导入包路径> <版本> [// indirect]`
    + `replace` 替换依赖项模块 `$module => $newmodule`。`$newmodule`可以是本地磁盘的相对/绝对/网络路径
- `go.mod`被创建后，它的内容将会被`go toolchain`全面掌控。`go toolchain`会在各类命令执行时维护、修改`go.mod`文件。命令如`get`, `build`, `mod`



## pprof

### 帮助链接
- [pprof命令详解](https://www.kancloud.cn/cattong/go_command_tutorial/261357)
- [pprof](https://www.cnblogs.com/qcrao-2018/p/11832732.html)


### 概要文件
- `概要文件`: 通过标准库的`runtime`和`runtime/pprof`中的程序能生成包含`实时性数据`的`概要文件`


#### CPU概要文件
- Go语言运行时系统会以`100Hz`的频率对`Goroutine堆栈上`的`程序计数器`进行取样，即每`10ms`取样一次
- `CPU主频`即CPU内核工作的时钟频率`CPU Clock Speed`(单位HZ)，`时钟频率`的倒数即为`时钟周期`，在一个`时钟周期`内，CPU执行一条`运算指令`。
- `pprof.StartCPUProfile(file)`开始记录`cpu使用情况`
- `pprof.stopCPUProfile()`停止记录`cpu使用情况`(取样的频率设置为0)

#### 内存概要文件
- Go语言运行时系统会记录`用户程序运行期间`的所有`堆内存分配`。只要`堆内存`被分配`MemProfileRate`，分析器就会对其进行取样
- `pprof.WriteHeapProfile(file)`将`堆内存分配情况`写入文件
- 分析器的取样间隔`runtime.MemProfileRate byte`: 分析器会在每分配指定的`字节数量`后对`内存使用情况`进行取样
    + 默认是`512*1024`即`512K`个字节
    + 我们将`MemProfileRate`赋值为`0`表示取消取样

#### 程序阻塞概要文件
- `程序阻塞概要文件`用于保存用户程序中的`Goroutine阻塞事件`的记录
- 设置分析器的取样间隔`runtime.SetBlockProfileRate(*rate)`: 每发生几次Goroutine阻塞事件时对这些事件进行取样。
    + 默认为1，即每次阻塞都会取样
    + 设置为0或负数，表示取消取样


#### 其他概要文件
- 通过`pprof.Lookup("block").WriteTo(file, 0)`将保存在运行时内存中的`内存使用情况`记录取出，并在记录的实例上调用`WriteTo`方法将记录写入到文件中。

|     名称     |           说明          |       默认取样频率      |
|--------------|-------------------------|-------------------------|
| goroutine    | 活跃goroutine的信息记录 | 获取时取样一次          |
| threadcreate | 系统线程创建情况的记录  | 获取时取样一次          |
| heap         | 堆内存情况分配的记录    | 默认每分配512KB取样一次 |
| block        | goroutine阻塞时间的记录 | 默认没次阻塞时取样      |
