package slog

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// fixedTime 用于固定日志里的时间，方便断言；真实日志里是调用时刻。
var fixedTime = time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

// Replacer 是 HandlerOptions.ReplaceAttr 的实现. groups 表示函数调用链中的组名
// 根属性的 groups 是 []; f21 的 groups 是 ["f2"];
func Replacer(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && a.Key == slog.TimeKey {
		// 内置 time, level, msg, source 四个属性, 将 time 值固定
		return slog.String("ts", fixedTime.Format(time.RFC3339))
	}

	// 处理 source 属性，只保留文件名
	if a.Key == slog.SourceKey {
		// 默认 source 会带上本机绝对路径和行号。
		if src, ok := a.Value.Any().(*slog.Source); ok {
			return slog.String(slog.SourceKey, filepath.Base(src.File))
		}
	}

	switch a.Key {
	case "password": // 脱敏
		return slog.String("password", "[REDACTED]")
	case "secret": // 返回零值 Attr 表示删除
		return slog.Attr{}
	}
	return a
}

// account 屏蔽敏感词
type account struct {
	name  string
	token string
}

// LogValue 实现 slog.LogValuer 接口，用于自定义日志输出。
func (a account) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", a.name),
		slog.String("token", "****"),
	)
}

// 结构化日志 slog 的基本用法
func Example_basic() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelDebug, // slog 的默认级别是 Info
		AddSource:   true,            // 记录调用日志的文件:行号（有性能开销，按需开启）
		ReplaceAttr: Replacer,        // 在每条属性输出前被调用. 可以修改或删除属性。
	}))

	// f1=1 是范型, 若要强类型, 可以使用 slog.Int("f1", 1)
	logger.Info("test info slog", "f1", 1)

	// LogAttrs 使用强类型属性记录日志。 两种输出完全一致
	logger.LogAttrs(context.Background(), slog.LevelInfo, "test info slog", slog.Int("f1", 1))

	logger.Warn("typed",
		slog.Float64("pi", 3.14),
		slog.Time("at", fixedTime),
		slog.Duration("cost", 5*time.Millisecond),
		slog.Group("f2", slog.String("f21", "v21")),
		slog.Any("err", fmt.Errorf("boom")),
	)

	// 实现 slog.LogValuer 接口的类型会自动调用 LogValue 方法。
	logger.Info("login", "account", account{name: "alice", token: "secret-token"})

	// Output:
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","source":"example_test.go","msg":"test info slog","f1":1}
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","source":"example_test.go","msg":"test info slog","f1":1}
	// {"ts":"2024-01-02T03:04:05Z","level":"WARN","source":"example_test.go","msg":"typed","pi":3.14,"at":"2024-01-02T03:04:05Z","cost":5000000,"f2":{"f21":"v21"},"err":"boom"}
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","source":"example_test.go","msg":"login","account":{"name":"alice","token":"****"}}
}

// 子 Logger 和日志分组
func Example_child() {
	parent := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: Replacer}))

	// With / WithGroup 返回新的 Logger，不会修改原 Logger。
	child := parent.WithGroup("http").With("req_id", "r-1")
	child.Info("child info")
	parent.Info("parent no req_id")

	// WithGroup 返回子 Logger，之后 With 和日志调用里的属性都会被加上组名前缀。
	parent.With("service", "order").WithGroup("http").Warn("request", "status", 200)
	// 反过来：先 With 的属性在组外，先 WithGroup 的属性在组内。
	parent.WithGroup("db").With("driver", "mysql").Warn("query", "ms", 3)

	// Output:
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","msg":"child info","http":{"req_id":"r-1"}}
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","msg":"parent no req_id"}
	// {"ts":"2024-01-02T03:04:05Z","level":"WARN","msg":"request","service":"order","http":{"status":200}}
	// {"ts":"2024-01-02T03:04:05Z","level":"WARN","msg":"query","db":{"driver":"mysql","ms":3}}
}

// 标准库 log 桥接; 旧代码改造 log -> slog
func Example_logBridge() {
	var lv slog.LevelVar // 零值 = INFO
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       &lv, // 用 Leveler（LevelVar 实现了它）动态取级别
		ReplaceAttr: Replacer,
	})
	logger := slog.New(handler)

	// slog.SetDefault，将标准库log转为slog，交给新的 Handler 输出
	slog.SetDefault(logger)
	log.Println("info bridge")

	// SetLogLoggerLevel 设置标准库 log 桥接的日志级别
	slog.SetLogLoggerLevel(slog.LevelDebug)
	log.Println("debug bridge") // 默认标准库 log 桥接的日志级别是 Info, debug会被丢弃
	lv.Set(slog.LevelDebug)     // 级别改成 DEBUG, debug 日志会被输出

	// logger.Enabled 判断指定级别是否会被输出
	if logger.Enabled(context.Background(), slog.LevelDebug) {
		// 典型用法：只有确定要打日志时，才去拼装耗时/昂贵的参数。
		log.Println("debug bridge")
	}

	// Output:
	// {"ts":"2024-01-02T03:04:05Z","level":"INFO","msg":"info bridge"}
	// {"ts":"2024-01-02T03:04:05Z","level":"DEBUG","msg":"debug bridge"}
}

// 把 Handler 包装成 *log.Logger：slog -> log 方向
func Example_logLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: Replacer})

	// 适合改造只接受 *log.Logger 的旧代码/第三方库。
	stdLogger := slog.NewLogLogger(handler, slog.LevelWarn)
	stdLogger.Print("from log.Logger")

	// Output:
	// {"ts":"2024-01-02T03:04:05Z","level":"WARN","msg":"from log.Logger"}
}
