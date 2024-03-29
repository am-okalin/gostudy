## concurrent
### 进程与线程
- `进程`是程序的运行时。`进程`是`资源管理`(内存、IO设备)的基本单位，系统会在创建/关闭`进程`时分配/回收资源(内存、IO设备)。因此进程开销大
- `线程`是`程序执行`(CPU)的基本单位，可访问`进程`的资源，可被视为`进程`中的代码执行的流程
- 每个`进程`的第一个`线程`都会随着该`进程`的启动而被创建，它们被称为其所属进程的`主线程`。
- 除了`主线程`外，其他的`线程`都是由`进程`中已`存在的线程`创建出来的。主线程之外的其他线程都只能由代码显式地创建和销毁。这需要我们在编写程序的时候进行手动控制

### 线程与协程`coroutines`
- 线程是同步机制、并行的、抢占式、内核态的
- 协程是异步机制、并发的、协作式、用户态的、更轻量的

### GO中的并发
- 在Go程序当中，`runtime`(运行时)会帮助我们自动地创建和销毁`系统级线程`。
- 而对应的用户级线程指的是架设在系统级线程之上的，用户级线程的创建、销毁、调度、状态变更以及其中的代码和数据都完全需要用户去实现和处理。

### GPM 并发编程模型
- `machine` 系统级线程(GO的运行时`runtime`管理生命周期)
- `goroutine` 可理解为用户级线程(开发者管理生命周期)
- `processor` 调度器，用于调度G与M对接。默认P的数量等于cpu核心数
    + 在等待I/O或者锁解除时P会分离对应的G与M
    + 在G需要恢复时，P会寻找M将两者对接
    + M不够用时P会向操作系统申请新的M
- 执行go语句时会先从`某个存放空闲的G队列`中获取一个G，找不到G时才会创建G；然后用G去包装go语言的函数把它追加到`某个存放可运行的G队列`中。因此
    + 已存在的G总会被优先复用
    + go函数的执行时间总是会明显滞后于它所属的go语句的执行时间

### goroutines
```go
f()    // call f(); wait for it to return
go f() // create a new goroutine that calls f(); don't wait
```

- `goroutines`是GO中的`并发执行单元`，可把它看作是一个`用户级的线程`。每个`goroutines`都有一唯一标识。
- 当一个程序启动时，其主函数在一个单独的`goroutine`中运行，我们叫它`main goroutine`，`main goroutine`结束时，会打断其所有的`子goroutine`

### channel
- `Communicating Sequential Process`简称`CSP`，中文直译`通信顺序进程`。CSP 允许使用`进程组件`来描述系统，它们独立运行，并且只通过`消息传递`的方式通信。
    + Don’t communicate by sharing memory; share memory by communicating.
    + 不要通过共享数据来通讯，要以通讯的方式共享数据。
- go通过引入`channel`实现`CSP`思想，其主要应用场景有
    + `数据交流`: 当作并发的`buffer`或`queue`。将`goroutine`当作生产者`Producer`和消费者`Consumer`
    + `数据传递`: `goroutine`间的数据传递，相当于把数据的拥有权 (引用) 托付出去。
    + `信号通知`: 一个`goroutine`可以将信号`closing、closed、data ready等` 传递给另一个或者另一组`goroutine`
    + `任务编排`: 让一组`goroutine`按顺序并发或者串行的执行
    + `锁`: 利用`channel`实现锁机制
- 使用反射操作channel

### channel下的任务编排
- Or-Done 任意一个`<-inchan`完成后，就关闭`<-outchan`
- 扇入模式 每个`inchan`都写入数据，一个`outchan`输出数据
- 扇出模式
- Stream
- map-reduce


### 名词解释
- `临界区` 共享资源，可以是IO操作，数据结构，变量
- `data race` 数据竞争/竞态条件，多线程对`临界区`的并发读写
- `Mutex` 排他/互斥锁，通过阻塞的方式解决`data race`的问题。
- `重入锁` 同个线程可对`临界区`多次加锁解锁,可递归调用，`mutex`是`不可重入锁`
- `死锁` 多个线程对多个`临界区`的相互持有与等待。
    + 如go1锁住v1请求v2,go2锁住v2请求v1就会导致死锁

### 并发原语的应用范围
- 共享资源的并发访问使用传统并发原语；
- 复杂的任务编排和消息传递使用 Channel；
- 消息通知机制使用 Channel，除非只想 signal 一个 goroutine，才使用 Cond；
- 简单等待所有任务的完成用 WaitGroup，也有 Channel 的推崇者用 Channel，都可以；
- 需要和 Select 语句结合，使用 Channel；
- 需要和超时配合时，使用 Channel 和 Context。

### 工具
- [race-detector](https://go.dev/blog/race-detector) 可做`data race`检测，只要在`test/run/build/install`时加上`-race`参数就可以检测了
- [go vet](#go-tool-vet) 可检查`死锁`([chekdead方法](https://go.dev/src/runtime/proc.go?h=checkdead#L4935))，它是通过[copylock分析器](https://github.com/golang/tools/blob/master/go/analysis/passes/copylock/copylock.go)静态分析实现的。

## context
### 基本使用方法
```golang
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

- `Deadline`: `deadline`指`ctx`被取消的截止日期。没设置截至日期则`ok=false`。后续调用返回结果相同。
- `cancel`、`timeout`、`deadline`都会导致`ctx`被`cancel`
- `Done` 根据`ctx`的类不同，有不同的值或处理方式
    + ctx不为`cancelCtx`: `Done()`返回`nil`
    + ctx为`cancelCtx`且未被`cancel()`: `Done()`返回阻塞的`chan`
    + ctx为`cancelCtx`且被`cancel()`: `Done()`返回被关闭的`chan`。`Err`方法返回`close`的原因。
- `Value` 返回此 `ctx` 中和指定的`key`相关联的`value`。
- `context.Background()`返回一个非`nil`的、空的`Context`，没有任何值，不会被 cancel，不会超时，没有截止日期。一般用在主函数、初始化、测试以及创建根`Context`的时候。
- `context.TODO()`底层实现与`context.Background()`一样。当你不清楚是否该用 Context，或者目前还不知道要传递一些什么上下文信息的时候，就可以使用这个方法。

### context使用规范
- 把`Context`作为方法的第一个参数
- 不使用`nil`作为`Context`的参数值
- `Context`只用来临时做函数之间的上下文透传，不能持久化`Context`或者把`Context`长久保存。把`Context`持久化到数据库、本地文件或者全局变量、缓存中都是错误的用法。
- `key`的类型不应该是`字符串类型`或者其它`内建类型`，否则容易在包之间使用`Context`时候产生冲突。使用`WithValue()`时，`key`的类型应该是`自定义类型`(非必须)
- 使用`struct{}`作为底层类型定义`key`的类型。使用`接口`或者`指针`作为底层类型定义`exported key`的静态类型。这样可以尽量减少内存分配

### context使用场景
- 上下文传递`request-scoped`: 如处理http请求，处理链路上的数据传递
- 控制子goroutine运行
- 超时控制的方法调用
- 可取消的方法调用

### WithValue
```golang
type valueCtx struct {
    Context
    key, val interface{}
}
```

- `WithValue`基于`Context intesrface`创建了`valueCtx`的实例。它持有一个KV键值对(用于传递上下文)。
- `valueCtx`覆盖了`Value`方法使用`链式查找`,它优先从自己的存储中检查这个`key`，不存在则从`Context intesrface`中继续检查，若`Context intesrface`也是`valueCtx`，则重复此过程(装饰器模式)

### WithCancel
- `WithCancel(parent Context) (ctx Context, cancel CancelFunc)` 返回的`ctx`为`cancelCtx`类型，会新建`ctx.Done`对象，用于取消长时间的任务。触发如下情况时`ctx.Done`就会被`close`
    + `parent.Done`被`close`时
    + 执行`WithCancel`返回的`cancel()`方法时
- `WithCancel`执行时会调用`propagateCancel`方法，此方法会顺着`parent`向上查找到一个`cancelCtx`或`nil`
    + 找到`nil`(根ctx)，就会新起一个`goroutine`，用于监听`parent.Done`是否已关闭。
    + 找到`cancelCtx`就把`当前ctx`加入到这个`cancelCtx`的`child`，以便这个`cancelCtx`被取消的时候通知`当前ctx`
- `cancel()`是向下传递的，`子孙ctx`会被`cancel()`，但`祖先ctx`不会被`cancel()`
- 注: 只要任务完成(正常或异常结束)，就要调用`cancel`。这样才可以释放`ctx`资源(通知` child`处理`cancel`；从`parent`的`child`中移除自己；甚至释放相关的`goroutine`)

### WithDeadline
- `WithDeadline(parent Context, d time.Time) (ctx Context, cancel CancelFunc)` d为截止时间，ctx可能为`cancelCtx`或`timerCtx`。
    + `d`晚于`parent的截止时间`则以后者为准，返回`cancelCtx`类型实例
    + 若`当前时间`超过了`截止时间`, 则返回已`cancel`的`timerCtx`。否则启动一个`timer`,到`截止时间`取消这个`timerCtx`
- `timerCtx.Done()`被`close`，有以下时间触发
    + 截止时间到了
    + `cancel` 函数被调用
    + `parent.Done`被`close`
- `timerCtx`也实现了`canceler`接口。不可依赖截止时间被动取消，`cancel`一定要尽早调用，这样才能尽早释放资源。

### WithTimeout
```golang
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    // 当前时间+timeout就是deadline
    return WithDeadline(parent, time.Now().Add(timeout))
}
```
WithTimeout 其实是和 WithDeadline 一样，只不过一个参数是超时时间，一个参数是截止时间。超时时间加上当前时间，其实就是截止时间
