## references
- [pkg.go.dev](https://pkg.go.dev/)
- [goproxy](https://goproxy.io/zh/)

### 官方文档
- [go_download](https://go.dev/dl/)
- [go_cmd](https://golang.google.cn/cmd/)
- [go_cmd_zh](https://www.kancloud.cn/cattong/go_command_tutorial/261347)
- [go_env](https://golang.google.cn/cmd/go/#hdr-Environment_variables)
- [go_mod](https://golang.org/ref/mod)
- [go_effective](https://go-zh.org/doc/effective_go.html)
- [go官方包](https://pkg.go.dev/std)

### 比较好的入门教程
- [go语言圣经](https://books.studygolang.com/gopl-zh/)
- [go语言高级编程](https://chai2010.cn/advanced-go-programming-book/)

### 知名的开源项目
- https://github.com/avelino/awesome-go
- https://github.com/golang-standards/project-layout


### 一些经验
- [go语言圣经](https://books.studygolang.com/gopl-zh/)是广受认可的入门书籍，另外官方教程也有中文版[go_effective](https://go-zh.org/doc/effective_go.html)。大致看下数据类型、控制结构、函数与方法等内容就可以进入开发了，遇到遗忘或不会时直接通过这两个文档进行查询。编写代码时通过测试文件快速调试。
- GO可开发的`应用`范围很广(web，游戏，网络编程，cli工具)，开发者根据自己的需求编写或自由选择`代码包`组装实现`应用`，定义项目的目录结构即`layout`
	+ GO社区对`layout`有一定的共识，参考[https://github.com/golang-standards/project-layout/]
- [https://github.com/avelino/awesome-go]() 收集了市面上大部分`轮子`，这个项目很重要，通常我们为完成某一应用要做技术选型时可在这个仓库中选择一些star较高的项目实现。
	+ 如要实现web应用可在`Web Frameworks`中找，要实现命令行工具可在`Command Line`中找，多线程可在`Goroutines`中找
- 可通过[pkg.go.dev](https://pkg.go.dev/)查看代码包的说明文档。官方包的使用可直接查看[pkg.go.dev/std](https://pkg.go.dev/std)。
	+ 但我认为通过`代码包`的`测试文件`/`样例文件`了解使用方式最方便
- 在go中学习接口 > 学习实现。应当面向接口测试。
