package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	rdb := redis.NewClient(redisOptions())

	result, err := rdb.Ping(rdb.Context()).Result()
	t.Log(result, err)
}

func TestExists(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(redisOptions())
	k1 := "k1"

	i64, err := rdb.ZAdd(ctx, k1, &redis.Z{Score: 0, Member: "member_1"}).Result()
	t.Log("添加总数 ", i64, err)

	i64, err = rdb.Exists(ctx, k1).Result()
	t.Log(i64, err)

	i64, err = rdb.ZRem(ctx, k1, "member_1").Result()
	t.Log("移除member_1 影响行数:", i64, err)

	i64, err = rdb.Exists(ctx, k1).Result()
	t.Log(i64, err)

	i64, err = rdb.Del(ctx, k1).Result()
	t.Log("删除key 影响行数:", i64)

	i64, err = rdb.Exists(ctx, k1).Result()
	t.Log(i64, err)

}

func TestString(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(redisOptions())

	err := rdb.Set(ctx, "key01", "value01", 0).Err()
	t.Log(err)
	flag, err := rdb.SetNX(ctx, "key01", "value02", 10*time.Second).Result()
	t.Log(flag, err)
	val, err := rdb.Get(ctx, "key01").Result()
	t.Log("key01", val)

	err = rdb.Set(ctx, "key02", 100, 0).Err()
	t.Log(err)
	ui64, err := rdb.Incr(ctx, "key02").Uint64()
	t.Log(ui64, err)
	i64, err := rdb.Get(ctx, "key02").Int64()
	t.Log("key02", i64)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(redisOptions())

	rdb.LPush(ctx, "list01", "val1", "val2")
	lRange := rdb.LRange(ctx, "list01", 0, -1)
	flag := Expect(lRange.Err()).NotTo(HaveOccurred())
	list01 := lRange.Val()
	t.Log(flag, list01)
}

func TestSortSet(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(redisOptions())
	zk1 := "zk1"

	var zl []*redis.Z
	for i := 0; i < 5; i++ {
		zl = append(zl, &redis.Z{
			Score:  float64(i),
			Member: fmt.Sprintf("member_%d", i),
		})
	}

	i64, err := rdb.ZCard(ctx, zk1).Result()
	t.Log("总数 ", i64, err)

	i64, err = rdb.ZAdd(ctx, zk1, zl...).Result()
	t.Log("添加总数 ", i64, err)

	i64, err = rdb.ZRem(ctx, zk1, "member_4").Result()
	t.Log("移除member_4 影响行数:", i64, err)

	i64, err = rdb.ZCard(ctx, zk1).Result()
	t.Log("总数 ", i64, err)

	rdb.ZIncrBy(ctx, zk1, 100, "member_1")
	f64, err := rdb.ZIncrBy(ctx, zk1, 100, "member_5").Result()
	t.Log(f64, err)

	f64, err = rdb.ZScore(ctx, zk1, "member_1").Result()
	t.Log("member_1的值 ", f64, err)

	i64, err = rdb.ZCount(ctx, zk1, "(1", "+inf").Result()
	t.Log("集合中大于1的数量", i64)

	sl, err := rdb.ZRevRange(ctx, zk1, 0, -2).Result()
	t.Log("从大到小 0至倒数第2个成员", sl, err)

	sl, err = rdb.ZRangeByScore(ctx, zk1, &redis.ZRangeBy{Min: "3", Max: "+inf"}).Result()
	t.Log("从小到大 分数大于等于3的成员", sl, err)

	zll, err := rdb.ZRangeWithScores(ctx, zk1, 0, 3).Result()
	t.Log("输出列表", zll, err)

	i64, err = rdb.ZRevRank(ctx, zk1, "member_11").Result()
	t.Log("从大到小 member_1的排名", i64, err)

	i64, err = rdb.Del(ctx, zk1).Result()
	t.Log("删除key 影响行数:", i64)
}
