## concurrent

### GPM 并发编程模型
- `machine` 系统级线程(GO的运行时`runtime`管理生命周期)
- `goroutine` 用户级线程(开发者管理生命周期)
- `processor` 调度器，用于调度G与M对接。通常P的数量等于cpu核心数
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

- `goroutines`是GO中的`并发执行单元`，可把它看作是一个`用户级的线程`(开发者管理生命周期)。与之相对的是`系统级线程`(GO的运行时`runtime`管理生命周期)。
- 当一个程序启动时，其主函数即在一个单独的`goroutine`中运行，我们叫它`main goroutine`，`main goroutine`结束时，会打断其所有的`子goroutine`
- 每个`goroutines`都有一唯一标识