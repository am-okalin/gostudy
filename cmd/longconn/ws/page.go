package main

import (
	"io"
	"net/http"
)

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, indexHTML)
}

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>WebSocket 聊天室</title>
<style>
  body { font: 14px/1.6 -apple-system, "PingFang SC", sans-serif; margin: 2rem; max-width: 720px; }
  #log { border: 1px solid #ddd; border-radius: 6px; padding: .75rem; height: 340px; overflow: auto; background: #fafafa; }
  .sys { color: #999; font-style: italic; }
  .row { display: flex; gap: .5rem; margin-top: .75rem; }
  input { flex: 1; padding: .4rem .6rem; }
  button { padding: .4rem 1rem; }
  .meta { color: #888; font-size: 12px; }
</style>
</head>
<body>
<h3>WebSocket 聊天室</h3>
<p class="meta">打开两个标签页就能互相对话。服务端每 30s ping 一次，浏览器会自动回 pong。</p>
<div id="log"></div>
<div class="row">
  <input id="text" placeholder="说点什么，回车发送">
  <button id="send">发送</button>
</div>
<script>
const name = prompt('你的名字？', '游客' + Math.floor(Math.random() * 1000)) || '匿名';
const log = document.getElementById('log');
const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
const ws = new WebSocket(proto + location.host + '/ws?name=' + encodeURIComponent(name));

ws.onopen = () => add('已连接', true);
ws.onclose = () => add('连接已关闭，刷新页面重连', true);
ws.onmessage = (e) => {
  const m = JSON.parse(e.data);
  add(m.system ? m.text : m.time + ' ' + m.name + ': ' + m.text, m.system);
};

function add(text, sys) {
  const div = document.createElement('div');
  if (sys) div.className = 'sys';
  div.textContent = text;
  log.appendChild(div);
  log.scrollTop = log.scrollHeight;
}

function send() {
  const el = document.getElementById('text');
  const text = el.value.trim();
  if (!text || ws.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({ text: text }));
  el.value = '';
}

document.getElementById('send').onclick = send;
document.getElementById('text').addEventListener('keydown', (e) => {
  if (e.key === 'Enter') send();
});
</script>
</body>
</html>`
