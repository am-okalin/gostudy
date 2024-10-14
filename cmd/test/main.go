package main

import (
	"fmt"
	"runtime"
	"time"
)

type ClickEvent struct {
	// EventId 业务端定义的uuid
	EventId string `json:"event_id,omitempty"`
	// Hostname req.Header.Hostname
	Hostname string `json:"hostname,omitempty"`
	// Path req.Route.Path
	Path string `json:"path,omitempty"`
	// ClientId 请求客户端标识 从配置文件中读取<cmp_shortlink>
	ClientId string `json:"client_id,omitempty"`
	// Referer req.Header.Referer
	Referer string `json:"referer,omitempty"`
	// IP req.Header.xFormatfor
	Ip string `json:"ip,omitempty"`
	// UserAgent req.Header.xFormatfor
	UserAgent string `json:"user_agent,omitempty"`
	// Language Accept-Language 没有就不传
	Language string `json:"language,omitempty"`
	// ClickedAt 创建时间, iso8601
	ClickedAt time.Time `json:"clicked_at,omitempty"`
	// Metadata 元数据, 本项目无内容
	Metadata interface{} `json:"metadata,omitempty"`
}

func main() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for i := 0; i < 1000; i++ {
		go func() {
			do := &ClickEvent{
				Hostname:  "localhost:8000",
				Path:      "/9f08f488",
				UserAgent: "Apifox/1.0.0 (https://apifox.com)",
				ClickedAt: time.Now(),
			}
			time.Sleep(time.Minute)
			fmt.Println(do)
		}()
	}

	for true {
		time.Sleep(time.Second)
		fmt.Println(runtime.NumGoroutine(), memStats.Alloc)
	}
}
