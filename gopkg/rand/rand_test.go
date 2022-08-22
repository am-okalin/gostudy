package rand

import (
	"crypto/rand"
	"math/big"
	mr "math/rand"
	"testing"
	"time"
)

// 伪随机, 需要独立的seed
func Test1(t *testing.T) {
	mr.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		println(mr.Intn(100))
	}
}

// 真随机, 性能慢
func Test2(t *testing.T) {
	for i := 0; i < 4; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		println(n.Int64())
	}
}
