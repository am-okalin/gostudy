package prod

import (
	"fmt"
	"kit/tableconv"
	"strings"
	"testing"
)

const (
	HeyupTryout = "Heyup Tryout Product - Tryout Product.csv"
	TryoutReal  = "tryout_real.csv"
	TryoutFake  = "tryout_fake.csv"
)

// 数据过滤 & 生成主键
func TestTryoutFilter(t *testing.T) {
	//加载已有csv
	from, err := tableconv.Csv2Table(HeyupTryout, ',')
	if err != nil {
		t.Error(err)
	}

	//用map[id=>arr]过滤重复
	m := make(map[string][]string)
	for _, row := range from {
		if row[0] == "Shop Name" {
			continue
		}
		row[3] = strings.ToLower(row[3])
		id := fmt.Sprintf("%s_%s_%s", row[3], row[0], row[2])
		m[id] = []string{id, row[3], row[0], row[2]}
	}

	//导入至新的table过程中做数据转换
	to := make([][]string, 1, len(m))
	to[0] = []string{"id", "email", "shop_name", "product_id"}
	for _, arr := range m {
		to = append(to, arr)
	}

	//导出新的csv
	err = tableconv.ToCsv(to, TryoutReal)
	if err != nil {
		t.Error(err)
	}
	t.Log("done")
}

// 生成假数据
func TestTryoutFake(t *testing.T) {
	//初始数据
	m := map[string]int{
		"7972537270510": 813,
		"7962577436910": 322,
		"7945430139118": 2197,
		"7945442787566": 769,
		"7945430761710": 1558,
		"7945431580910": 647,
		"7945430925550": 360,
		"7962490142958": 338,
		"7945431810286": 4498,
		"7945431089390": 375,
		"7945449406702": 419,
		"7945432006894": 217,
	}

	shopName := "heyup-web"

	//生成新的table
	var length int
	for _, count := range m {
		length += count
	}
	to := make([][]string, length+1)
	index := 0
	to[index] = []string{"id", "email", "shop_name", "product_id"}
	for productId, count := range m {
		for i := 0; i < count; i++ {
			index++
			email := fmt.Sprintf("%d@test.com", index)
			id := fmt.Sprintf("%s_%s_%s", email, shopName, productId)
			to[index] = []string{id, email, shopName, productId}
		}
	}

	//导出新的csv
	err := tableconv.ToCsv(to, TryoutFake)
	if err != nil {
		t.Error(err)
	}
	t.Log("done")
}
