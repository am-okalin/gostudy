// 裸 UDP 示例：没有连接、没有重传、没有顺序，全靠应用自己兜。
//
// 这个 demo 想让你直观看到三件事：
//  1. 发出去 ≠ 收到（服务端 -drop 20% 模拟丢包）；
//  2. 收不到也不会报错，UDP 不会告诉你"对方没回"；
//  3. 想可靠，就得自己编号 + 超时 + 重传——这正是 QUIC 存在的理由。
package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

// packet 是 UDP 的数据报格式（一条消息 = 一个 datagram，天然有边界，不用分帧）。
type packet struct {
	Seq        int    `json:"seq"`
	SentNano   int64  `json:"sent_nano"`
	ServerNano int64  `json:"server_nano,omitempty"`
	Text       string `json:"text,omitempty"`
}

func main() {
	addr := flag.String("addr", "127.0.0.1:9401", "监听地址；客户端模式下是服务端地址")
	asClient := flag.Bool("client", false, "以客户端身份运行")
	drop := flag.Float64("drop", 0.2, "服务端模拟丢包率（0~1）")
	count := flag.Int("count", 10, "客户端发送的包数")
	flag.Parse()

	if *asClient {
		if err := runClient(*addr, *count); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := runServer(*addr, *drop); err != nil {
		log.Fatal(err)
	}
}

func runServer(addr string, drop float64) error {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer pc.Close()

	log.Printf("UDP 服务启动: %s（模拟丢包率 %.0f%%）", pc.LocalAddr(), drop*100)

	buf := make([]byte, 64<<10)
	for {
		n, raddr, err := pc.ReadFrom(buf)
		if err != nil {
			return err
		}
		if drop > 0 && rand.Float64() < drop {
			continue // 故意不回包，模拟链路丢包
		}

		var m packet
		if err := json.Unmarshal(buf[:n], &m); err != nil {
			continue
		}
		m.ServerNano = time.Now().UnixNano()

		out, err := json.Marshal(m)
		if err != nil {
			continue
		}
		// UDP 的 WriteTo 不会因为对方不在而报错，发出去就"发出去"了
		pc.WriteTo(out, raddr)
	}
}

func runClient(addr string, count int) error {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	sent := make(map[int]time.Time, count)
	got := make(map[int]bool, count)

	send := func(seq int) {
		m := packet{Seq: seq, SentNano: time.Now().UnixNano(), Text: "hello"}
		b, _ := json.Marshal(m)
		conn.Write(b) // 忽略错误也没关系：UDP 本来就是"尽力而为"
		sent[seq] = time.Now()
	}

	// 1) 一口气把包全发出去
	for i := 0; i < count; i++ {
		send(i)
	}
	log.Printf("已发送 %d 个包，等 1.5s 收回应...", count)

	// 2) 收回应，直到超时
	collect := func(d time.Duration) int {
		buf := make([]byte, 4096)
		deadline := time.Now().Add(d)
		newly := 0
		for {
			conn.SetReadDeadline(deadline)
			n, err := conn.Read(buf)
			if err != nil {
				return newly // 超时 or 其他错误：这一轮结束
			}
			var m packet
			if err := json.Unmarshal(buf[:n], &m); err != nil || got[m.Seq] {
				continue
			}
			got[m.Seq] = true
			newly++
			log.Printf("  收到 seq=%-2d rtt=%.2fms", m.Seq, float64(time.Since(sent[m.Seq]).Microseconds())/1000)
		}
	}

	first := collect(1500 * time.Millisecond)
	total := first
	log.Printf("第一轮收到 %d/%d", first, count)

	// 3) 只重传没收到的：这就是"应用层可靠性"，QUIC 内部做的也是这件事
	var missing []int
	for i := 0; i < count; i++ {
		if !got[i] {
			missing = append(missing, i)
		}
	}
	if len(missing) > 0 {
		log.Printf("重传 %v ...", missing)
		for _, seq := range missing {
			send(seq)
		}
		total += collect(1500 * time.Millisecond)
	}

	log.Printf("最终收到 %d/%d", total, count)
	if total < count {
		log.Printf("仍有 %d 个包丢失——UDP 不会替你重试，这就是成本", count-total)
	}
	return nil
}
