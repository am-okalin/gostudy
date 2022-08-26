## todo
- map、结构体、方法、json转换
- goruntine 和 chan 的基本使用(优雅退出)
- context详解
- 从接口中往文件中读写内容
- 黑名单用户处理
- 错误处理的几种方式
- 如何请求外部接口

## 01 作用域、slice、func、type
- 上节内容: 包管理、简单测试、函数、引用类型、切片
- 什么情况下变量是可被其他包引用的
- 数组与切片的区别是什么? 可以类比成什么?(滑块[https://www.html5tricks.com/demo/vuejs-slider-pips/index.html])
- 变量的声明方式有哪些，有什么区别

## 02 map、struct
### map
- 切片->索引数组， 关联数组->键值对->字典/映射,在go中用map表示(如目录)
- GO是静态语言，每个类型必须声明，如何声明map，map是否要像slice一样初始化，不初始化的后果是什么
- map是引用类型，所以不初始化的赋值会报错，所以必须深拷贝

### struct
- 起源于C，C中的`struct`用于表示`复合类型`，在C中常用来组成`链表`，`队列`，`树`，`图`等数据结构。或根据特定的需求由开发者自定义的数据结构。
- 其实go中的各个复杂类型都是由`strcut`构成的(go在1.5时就实现了自举)
- 结构体的应用: 启一个`http`服务，关于`struct`和`json包`的简单应用
- 以上是struct表示具体类型时的应用，go作为21世纪C语言也实现了面向对象，即用struct表示抽象类型。->人类，汽车。struct拥有方法，可以类比对象的方法。
- 常会将方法看作对对象的操作或查询，指针方法和非指针方法有什么区别？
- struct的组合语法糖，可再次应用json包
- 字符操作相关类型: `strings` `[]rune` `[]byte`的应用与区别
	+ https://segmentfault.com/q/1010000009652523

## 03 面向接口编程
### interface
- 接口的定义，实现，测试->目前layout的应用
- 接口组合与面向接口编程(go极其鼓励):fmt.String,io.Reader,io.Writer,io.Closer,sort包,json包,hash包	

## 0905作业(3选2)
- 通过json字符串，获取结构体切片和map
- file包应用：新建文件，在文件中追加数据，从文件中读取数据
- http实现post接口获取外部数据并输出 

## 04 并发
- 面向接口编程的包讲解: io包相关，hash包相关，sort包相关
- GO中的并发讲解 -> sync包详解