package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Stream 处理 GET /events。
func (b *Broker) Stream(w http.ResponseWriter, r *http.Request) {
	// 关键：用 ResponseController，而不是 w.(http.Flusher)。
	// 中间件把 w 包了一层之后类型断言会失败，控制器会顺着 Unwrap() 找到真正的 writer。
	rc := http.NewResponseController(w)

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no") // 让 nginx 别缓冲

	w.WriteHeader(http.StatusOK)

	// 立刻 flush 一次，客户端才知道连接真的建起来了
	if err := rc.Flush(); err != nil {
		return
	}

	// retry: 断线后浏览器多久重连（毫秒）。空行表示一条消息结束。
	fmt.Fprint(w, "retry: 3000\n\n")
	fmt.Fprintf(w, "event: hello\ndata: 连接成功，等你收消息\n\n")
	if err := rc.Flush(); err != nil {
		return
	}

	events, unsub := b.Subscribe()
	defer unsub()

	// 心跳：15 秒一个注释行，防止中间的代理/负载均衡把空闲连接掐掉。
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			// 客户端断开（关页面、网络断），ctx 会被取消
			return

		case ev := <-events:
			if !writeEvent(w, ev) {
				return
			}
			if err := rc.Flush(); err != nil {
				return
			}

		case <-heartbeat.C:
			// 以 ":" 开头的行是注释，客户端会忽略，只用来保活
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return
			}
			if err := rc.Flush(); err != nil {
				return
			}
		}
	}
}

// writeEvent 按 SSE 格式写一条事件。
// 多行 Data 要给每一行都加 "data: " 前缀，否则只有第一行生效。
func writeEvent(w io.Writer, ev Event) bool {
	var b strings.Builder
	if ev.ID != 0 {
		fmt.Fprintf(&b, "id: %d\n", ev.ID)
	}
	if ev.Type != "" {
		fmt.Fprintf(&b, "event: %s\n", ev.Type)
	}
	for _, line := range strings.Split(ev.Data, "\n") {
		b.WriteString("data: ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteByte('\n') // 消息结束

	_, err := io.WriteString(w, b.String())
	return err == nil
}

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>SSE 演示</title>
<style>
  body { font: 14px/1.6 -apple-system, "PingFang SC", sans-serif; margin: 2rem; max-width: 720px; }
  #log { border: 1px solid #ddd; border-radius: 6px; padding: .75rem; height: 320px; overflow: auto; background: #fafafa; }
  .row { display: flex; gap: .5rem; margin-top: .75rem; }
  input { flex: 1; padding: .4rem .6rem; }
  button { padding: .4rem 1rem; }
  .meta { color: #888; font-size: 12px; }
</style>
</head>
<body>
<h3>Server-Sent Events 演示</h3>
<p class="meta">服务端每 5 秒推一条；也可以在下面手动广播一条，所有打开的页面都会收到。</p>
<div id="log"></div>
<div class="row">
  <input id="text" value="hello from browser">
  <button id="send">广播</button>
</div>
<script>
const log = document.getElementById('log');
function add(text) {
  const div = document.createElement('div');
  div.textContent = new Date().toLocaleTimeString() + '  ' + text;
  log.appendChild(div);
  log.scrollTop = log.scrollHeight;
}

const es = new EventSource('/events');
es.addEventListener('open', () => add('[open] 连接已建立'));
es.addEventListener('hello', (e) => add('[hello] ' + e.data));
es.addEventListener('message', (e) => add('[message] id=' + e.lastEventId + ' ' + e.data));
es.addEventListener('error', () => add('[error] 连接断开，浏览器会在 retry 后自动重连'));

document.getElementById('send').onclick = async () => {
  const text = document.getElementById('text').value;
  const resp = await fetch('/publish', { method: 'POST', body: JSON.stringify({ text: text }) });
  add('[post] /publish -> ' + resp.status);
};
</script>
</body>
</html>`

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, indexHTML)
}
