## todo
- map、结构体、方法、json转换
- goruntine 和 chan 的基本使用(优雅退出)
- context详解
- 从接口中往文件中读写内容
- 黑名单用户处理
- 错误处理的几种方式
- 如何请求外部接口

## 作用域、slice、func、type
- 上节内容: 包管理、简单测试、函数、引用类型、切片
- 什么情况下变量是可被其他包引用的
- 数组与切片的区别是什么? 可以类比成什么?(滑块[https://www.html5tricks.com/demo/vuejs-slider-pips/index.html])
- 变量的声明方式有哪些，有什么区别

## map
- 切片->索引数组， 关联数组->键值对->字典/映射,在go中用map表示(如目录)
- GO是静态语言，每个类型必须声明，如何声明map，map是否要像slice一样初始化，不初始化的后果是什么
- map是引用类型，所以不初始化的赋值会报错，所以必须深拷贝

## struct
- 起源于C，C中的`struct`用于表示`复合类型`，在C中常用来组成`链表`，`队列`，`树`，`图`等数据结构。或根据特定的需求由开发者自定义的数据结构。
- 其实go中的各个复杂类型都是由`strcut`构成的(go在1.5时就实现了自举)
- 结构体的应用: 启一个`http`服务，关于`struct`和`json包`的简单应用
- 以上是struct表示具体类型时的应用，go作为21世纪C语言也实现了面向对象，即用struct表示抽象类型。->人类，汽车。struct拥有方法，可以类比对象的方法。
- 常会将方法看作对对象的操作或查询，指针方法和非指针方法有什么区别？
- struct的组合语法糖，可再次应用json包
- 字符操作相关类型: `strings` `[]rune` `[]byte`的应用与区别
	+ https://segmentfault.com/q/1010000009652523

## interface
- 接口的定义，实现，测试->目前layout的应用
- 接口组合与面向接口编程(go极其鼓励):fmt.String,io.Reader,io.Writer,io.Closer,sort包,json包,hash包	

## 并发
- 主要包含三个内容 关键字goruntine 类型channel 包sync

### goruntine
- 复习面向接口编程，回顾case02的作业
- error,panic,defer,recover职能&特性
- goruntine介绍，使用方式，关闭进程与协程的类比，并行的数量控制
- goruntine_test.go test1->test2->testPN
- stack_test 例子中展示了如何获取唯一标识。实际上是较为复杂的字符串截取->go官方包不希望开发者通过`goid`进行协程控制

### channel01
- 复习：父子goruntine的关系，panic对goruntine的影响，
	+ 子goruntine发生panic时会导致什么问题，如何解决
	+ 能否往子goruntine中传递引用类型，为什么
	+ 上节课案例中我们是如何阻塞goruntine的?
- 阻塞队列，并发队列的定义，阻塞队列在什么时候阻塞? ->(不)带缓存的channel
- 手写代码： go1中输入字符，go2输出字符

### channel02
- channel如何实现协程间通信，定义
- wg并发的简单使用


## 作业
### 0905作业(3选2)
- 通过json字符串，获取结构体切片和map
- file包应用：新建文件，在文件中追加数据，从文件中读取数据
- http实现post接口获取外部数据并输出 

## 学习方法论
- 我要做并发往文件中写数据，等到做的那一刻再学行不行？不行！
	+ 现学现做无法保证质量
	+ 如此学习方法太片面，复杂知识点不能举一反三(下次要往reqBody中并发写数据)
- 分享会目的：掌握各个知识点的职能与特性(适合解决什么问题)