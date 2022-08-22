// Server2 is a minimal "echo" and counter server.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/post", poster)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// poster 输入姓民年龄，返回请求信息
func poster(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.URL, r.Method)

	var p person
	buf, _ := ioutil.ReadAll(r.Body)
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
