[TOC]
## io
- 没有嵌入其他接口并只定义了一个方法的接口叫`简单接口`
- 有着众多的`扩展接口`和`实现类型`的`简单接口`叫`核心接口` 
- io包中共有11个`简单接口`，其中4个是`核心接口`。`简单接口`针对`读` `写` `关闭` `设置读写位置`等操作分为四类。 接口如下
- io.Reader
- io.Writer
- io.Closer
- io.Seeker
- io.ReaderFrom 从`ReaderFrom`的参数`r`中读取数据写入到其`实例`
- io.WriterTo 从其`实例`中读取数据写入到`WriterTo`的参数`w`
- io.ReaderAt 单纯的只读方法，在其实现方法中不会移动`已读计数`,并发安全
- io.WriterAt
- io.ByteReader 读取下一个`byte字节`
- io.ByteScanner 内嵌/组合了`io.ByteReader`增加了一个`读回退单个字节`的功能集
- io.ByteWriter
- io.RuneReader 读取下一个`unicode字符`
- io.RuneScanner 内嵌/组合了`io.RuneReader`增加了一个`读回退单个unicode字符`的功能集
- io.StringWriter

## strings
- `string`类型内部有个指针指向底层字节数组的头部，但它仍然是值类型。
    + 值是不可变的，只能裁剪(切片)、拼接(+号)后生成一个新的字符串
    + `string值`会在底层与它的所有副本共用同一个`字节数组`。由于`字节数组`不可变所以绝对安全
- `strings.Builder` 与`string`类型同样拥有高效利用内存的前提条件。同时`Builder`支持内容追加(拼接)或完全重置，但其中内容仍不可修改/减少
    + 在已被真正使用后就不可再被复制,(复制后的任何方法都会引发panic)
    + 由于其内容不是完全不可变的，所以需要自行解决操作`冲突`和`并发安全`问题
- 自动扩容：`Builder`的拼接方法`Write`、`WriteByte`、`WriteRune`、`WriteString`会自动在`内容容器(字节数组)`容量不够用时进行扩容
- `b.Grow(n int)`手动扩容: 当`剩余容量`小于`n`时生成`2×旧容量+n`的新容器，将旧容器的数据拷贝至新容器
- `b.Reset()`重置`Builder`值重会零值状态
- `strings.Reader` 通过`已读计数`(用于读取索引，回退，位置设定)实现了高效读取字符串。
    + `Reader`大部分读取方法(`ReadByte`,`ReadRune`)都会更新`已读计数`，但`ReadAt`除外
    + `Seek(offset int64, whence int)`方法重新定位`计数`
    + 通过`r.Size()-int64(r.Len())`计算出`已读计数`

## bytes
- `bytes`与`strings`的api非常相似,不过它面对的主要是`字节`和`字节切片`。`strings`包主要面向`Unicode字符`和`经过UTF-8编码的字符串`
- `bytes.Buffer`即缓冲区，是集读写于一身的数据类型。使用`字节切片`作为`内容容器`；同时内部维护了一个`已读计数`。`Buffer已读计数`前的内容几乎没有机会再次被读取。
- 由于`Buffer`的`Cap()`方法提供的是`内容容器`的`容量`而不是`长度`，因此无法计算出`已读计数`
- `Truncate(n int)`截断方法`n`表示截断时保留`未读部分`头部的多少字节，此时`内容容器新长度=已读计数+n`
- 扩容：方法会在必要时依据`已读计数`找到未读部分，把内容拷贝到扩容后内容容器的头部后将`已读计数`置为0
- `Buffer`内容的泄露：`Bytes()`和`Next()`方法返回切片的底层数组与`Buffer`的底层数组一致。此时外界可更改这个数组导致严重的数据安全问题，所以在传出切片时要做好隔离(如对它们做`深拷贝`再把副本传出去)


## net
- `IPC`时`Inter-Process Communication` 的简写，表示进程间通信。主要定义的是多个进程之间，相互通信的方法。主要包括`signal`,`pipe`,`socket`,`file lock`,`message queue`,`semaphore`。在众多的`IPC`方法中，`socket`是最通用和灵活的一种。
- 现存的主流操作系统大都对`IPC`提供了强有力的支持，尤其是`socket`。支持`socket` 的操作系统一般都会对外提供一套`API`。跑在它们之上的应用程序利用这套`API`，就可以与互联网上的任意一台计算机上的程序，甚至同一个程序中的其他线程进行通信。
- `Go`对`IPC`也提供了一定的支持。
    + 在`os`代码包和`os/signal`代码包中就有针对系统信号的`API`。
    + `os.Pipe()`可以创建`命名管道`，而`os/exec`代码包则对另一类`管道(匿名管道)`提供了支持
    + 对于`socket`，`Go`与之相应的程序实体都在其标准库的`net`代码包中。


## os
- `socket(domain, stype, proto)`
- `DGRAM` 数据报 有消息边界，无逻辑连接
- `STREAM` 数据流 无消息边界，有逻辑连接

| domain/通信域 | stype/类型 | proto/协议 |
|---------------|------------|------------|
| ipv4          | DGRAM      | udp        |
| ipv6          | STREAM     | tcp        |
| unix          | SEQPACKET  |            |
|               | RAW        |            |

### File
#### 文件的操作模式
- 针对File值的主要操作模式有os.O_RDONLY、os.O_WRONLY、os.O_RDWR。必须把这三个模式中的一个设定为此文件的操作模式。
- os.O_APPEND：当向文件中写入内容时，把新内容追加到现有内容的后边。
- os.O_CREATE：当给定路径上的文件不存在时，创建一个新文件。
- os.O_EXCL：需要与os.O_CREATE一同使用，表示在给定的路径上不能有已存在的文件。
- os.O_SYNC：在打开的文件之上实施同步 I/O。它会保证读写的内容总会与硬盘上的数据保持同步。
- os.O_TRUNC：如果文件已存在，并且是常规的文件，那么就先清空其中已经存在的任何内容。

## wire
- 依赖项注入是一种标准技术，通过显式地为组件提供其工作所需的所有依赖项，从而生成灵活且松散耦合的代码。
- 在GO中，`依赖注入`通常采用将依赖项传递给构造函数的方式。它有两个基础的概念`providers`提供者； `injectors`注入者
- `providers`是普通的Go函数，它们根据依赖关系提供`provide`值，这些依赖关系被简单地描述为函数的参数。下面是一些示例代码，定义了三个提供程序
```go
// NewUserStore is the same function we saw above; it is a provider for UserStore,with dependencies on *Config and *mysql.DB.
func NewUserStore(cfg *Config, db *mysql.DB) (*UserStore, error) {...}

// NewDefaultConfig is a provider for *Config, with no dependencies.
func NewDefaultConfig() *Config {...}

// NewDB is a provider for *mysql.DB based on some connection info.
func NewDB(info *ConnectionInfo) (*mysql.DB, error) {...}
```
- 通常共用的`providers`可以分组到`providerSets`中
```go
var UserStoreSet = wire.ProviderSet(NewUserStore, NewDefaultConfig)
```
- `injectors` 是按`依赖顺序`调用`providers`的函数。函数中你需要书写注入器的签名`signature`，包括
    + 必要的参数
    + 在函数中调用`wire.Build(providers, providerSets, ...)`
```go
func initUserStore() (*UserStore, error) {
    wire.Build(UserStoreSet, NewDB)
    return nil, nil  // These return values are ignored.
}
```
- 执行`go generate`生成依赖关系描述文件`wire_gen.go`。任何非`injectores`的声明都会被复制到生成的文件中。程序运行时不依赖于`wire`:所有编写的代码都是普通的Go代码。


