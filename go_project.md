## go project 工程化
### 规范设计
- 非编码类规范： 开源规范 文档规范 版本规范 Commit规范 发布规范
- 编码类规范： 目录规范 代码规范 接口规范 日志规范 错误码规范

#### 开源规范
- 发布可用的版本：确保每次发布都经过充分测试(详细的单元测试)，每个发布版本都是可用的。
- 合适的git工作流。遵循`Angular commit message`规范。提交记录中不出现内部IP、内部域名、密码、密钥等信息
- `LICENSE`开源软件协议： `GPL`,`MPL`,`LGPL`,`Apache`,`BSD`,`MIT`
- `CONTERIBUTING.md`： 用于说明如何给本项目贡献代码，包含详细贡献流程
- `CHANGELOG`目录： 用于存放版本变更历史
- 详细的文档说明
    + `Maikefile`： 对项目进行构建、测试、安装等操作
    + `README.md`： 包含`项目描述`,`依赖项`,`安装方法`,`使用方法`,`贡献方法`,`作者`,`遵循软件协议`
    + `docs`目录： 存放项目所有文档如 `安装文档`,`使用文档`,`开发文档`
    + `examples`目录： 存放示例代码

### error
- go中没有`try_catch`, 设计者认为将异常与控制结构混在一起会很容易使得代码变得混乱。
- go中使用 error错误类型，panice异常，recover捕获异常，defer退出时执行等关键字或类型对错误进行处理

#### 错误处理
- 错误应该只处理一次，做一个决定
- 降级过的错误不应该继续上抛
- 应当优雅减少`Error Handling`。类似`bufio.Scanner`的`Scan()`和`Err()`的联合处理，大大减少了`client`的调用

#### 错误的几种定义方式
- `Sentinel Error`特定的错误，类似于`io.Eof`，更底层的有`syscall.ENOENT`。通常是用于开发时的调试，而不是业务开发时的返回
- `Error Type` 实现了`error接口`的自定义类型。如`os.PathError`。可以用`类型断言`或`类型switch`的方式通过`自定义类型`的变量获取更多的上下文
- `Opaque errors` 不透明的错误处理： 只返回错误而不假设其内容，通过行为获取更详细的信息。
    + 如`net.Error`在外部判断时想要更详细的信息可通过`Timeout`等方法获取更详细的信息

#### `github.com/pkg/errors` 堆栈错误包
- 可用`errors.Wrap`包装错误，并记录错误堆栈
- 可用`errors.Cause`获取根因，进行断言。
- 在错误产生处返回包装错误，在程序顶部用`%+v`记录堆栈详情
- 仅在业务、项目中使用`pkg/errors`，库中不应该使用。

#### go error新特性
- `errors.Is`和`errors.As`参数都可以对实现了`uwrap()`接口的错误进行递归判断