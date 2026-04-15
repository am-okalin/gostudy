package rand

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	mr "math/rand"
	"strconv"
	"strings"
	"testing"
)

// 真随机, 性能慢
func Test2(t *testing.T) {
	for i := 0; i < 4; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		println(n.Int64())
	}
}

// 伪随机, 需要独立的seed
func Test1(t *testing.T) {
	//mr.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		println(mr.Intn(10))
	}
	println(mr.Intn(10))
}

func Test3(t *testing.T) {
	str := GenerateRandomSteps(10)
	t.Log(str)
}

func GenerateRandomSteps(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := big.NewInt(int64(len(charset)))

	steps := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLength)
		steps[i] = charset[randomIndex.Int64()]
	}

	return string(steps)
}

func PowInt36(n int) *big.Int {
	base := big.NewInt(36)
	exponent := big.NewInt(int64(n))
	return new(big.Int).Exp(base, exponent, nil)
}

func GenerateCodes(generateNum, codeLength int, m map[string]bool) ([]string, error) {
	// 校验最大生成上限
	if generateNum > 1000000 {
		return nil, errors.New("超过最大生成上限")
	}

	max := PowInt36(codeLength)
	if big.NewInt(int64(generateNum)).Cmp(max) > 0 {
		log.Fatal(generateNum, max)
		return nil, errors.New("generateNum 必须小于可生成的总数")
	}

	// 生成list
	list := make([]string, 0, generateNum)
	for len(list) < generateNum {
		random, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		key := strings.ToUpper(strconv.FormatInt(random.Int64(), 36))

		for len(key) < codeLength {
			key = "0" + key
		}

		if !m[key] {
			list = append(list, key)
			m[key] = true
		}
	}
	return list, nil
}

func Test4(t *testing.T) {
	m := map[string]bool{"AA": true}
	list, err := GenerateCodes(100000, 255, m)
	t.Log(len(list), err, list)
}
