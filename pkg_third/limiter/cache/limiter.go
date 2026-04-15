package cache

import (
	"github.com/go-pkgz/expirable-cache/v3"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

const (
	defaultTokenBucketTTL = time.Minute
	defaultTokenBucketMax = 1000
)

type OptionFunc func(*TBLimiter)

type TBLimiter struct {
	b              int        // 令牌桶大小
	r              rate.Limit // 速率
	mu             *sync.Mutex
	tokenBuckets   cache.Cache[string, *rate.Limiter]
	tokenBucketTTL time.Duration // 令牌桶过期时间
	tokenBucketMax int           // 令牌桶最大数量
}

func TokenBucketTTL(ttl time.Duration) OptionFunc {
	return func(l *TBLimiter) {
		l.tokenBucketTTL = ttl
	}
}

func TokenBucketMax(max int) OptionFunc {
	return func(l *TBLimiter) {
		l.tokenBucketMax = max
	}
}

// NewTBLimiter 构造函数
func NewTBLimiter(r rate.Limit, b int, opts ...OptionFunc) *TBLimiter {
	tbl := &TBLimiter{
		b:              b,
		r:              r,
		mu:             &sync.Mutex{},
		tokenBucketTTL: defaultTokenBucketTTL,
		tokenBucketMax: defaultTokenBucketMax,
	}

	// 加载选项
	for _, o := range opts {
		o(tbl)
	}

	// 初始化令牌桶缓存
	tbl.tokenBuckets = cache.NewCache[string, *rate.Limiter]().WithTTL(tbl.tokenBucketTTL)
	return tbl
}

// GetLimiter 增加令牌桶 懒汉模式
func (tbl *TBLimiter) GetLimiter(key string) *rate.Limiter {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	// LRU算法自动删除过期的令牌桶
	limiter, found := tbl.tokenBuckets.Get(key)
	if found {
		return limiter
	}

	limiter = rate.NewLimiter(tbl.r, tbl.b)
	tbl.tokenBuckets.Add(key, limiter)
	return limiter
}
