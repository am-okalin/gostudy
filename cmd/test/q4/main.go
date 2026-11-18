// 最小 Web 服务器 + 原子计数，用来验证三件事：
//  1. net/http 默认 ServeMux 的路由与最长前缀匹配；
//  2. 多个请求 goroutine 共享变量时，读-改-写必须加锁；
//  3. 请求体的 JSON 反序列化与响应的 JSON 序列化。
//
// 源自《Go 程序设计语言》1.7 节的 Server2，另外追加了 /post 练手 JSON。
// 注：若本机 Docker 已占用 8000（监听 *:8000），要用 127.0.0.1:8000 访问才打得到这里。
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// count 被每个请求 goroutine 共享。lock 保护的范围是 “读值 + 写值”，
// 去除 Lock/Unlock 后用 `go run -race` 发起并发请求会直接报 DATA RACE。
var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/post", poster)
	// ListenAndServe 阻塞，log.Fatal 负责在绑定失败等场景直接退出。
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
// "/" 注册的是子树：未注册更具体路径的请求都落到这里，因此只有它们被计数。
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++ // 只有自增，但 count++ 是读-改-写，并非原子操作
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
// 读 count 同样要加锁，否则可能读到写了一半的值；真实现里应取出快照后再解锁，
// 把 Fprintf 这类网络 I/O 留在锁内会拉长临界区（此处从简）。
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// poster 输入姓名年龄，返回请求信息。
// 注意：这里既不校验 r.Method，也不校验解码结果——非法 JSON 只打日志，
// 仍会以 200 + 结构体零值响应。生产代码应校验 method 并对解码失败返回 400。
func poster(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.URL, r.Method)

	var p person
	buf, _ := io.ReadAll(r.Body) // 演示从简：读 Body 的错误被忽略
	err := json.Unmarshal(buf, &p)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "my name is %s, my age is %d\n", p.Name, p.Age)

	gary := person{
		Name: "gary",
		Age:  25,
	}

	garyJson, err := json.Marshal(gary)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "the json is %v\n", string(garyJson))
}
