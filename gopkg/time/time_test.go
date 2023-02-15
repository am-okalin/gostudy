package time

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestLocation(t *testing.T) {
	//定义时区为东八区
	loc, err := time.LoadLocation("Asia/Shanghai")
	t.Log(loc, err)

	//自定义location GMT
	loc = time.FixedZone("GMT", 8*3600)
	t.Log(loc)

	//获取time.UTC
	loc, err = time.LoadLocation("UTC")
	t.Log(loc, err)

	t1 := time.Now()
	t.Log(t1)

	//转换为UTC时区
	t1 = t1.UTC()
	t.Log(t1)

	//转换为本地时区
	t2 := time.Now().Local()
	t.Log(t2)
	t3 := time.Now().In(loc)
	t.Log(t3)

	//通过指定时间信息及loc得到Time
	t4 := time.Date(2001, 2, 3, 4, 5, 6, 0, loc)
	t.Log(t4)

	//时间字符串 转 time.Time
	t5, err := time.ParseInLocation("2006-01-02 15:04:05", "2001-02-03 04:05:06", loc)
	t.Log(t5, err)

	//默认以 UTC+0 的时区编译为time.Time
	t6, err := time.Parse("2006-01-02 15:04:05", "2001-02-03 04:05:06")
	t.Log(t6, err)
}

func TestTime(t *testing.T) {
	var q time.Time
	qu := q.Unix()
	t.Log(qu)

	a := time.Now()
	t.Log(a)

	// 返回秒数
	b := a.Unix()
	b2 := a.UnixMilli() //13位毫秒
	b3 := a.UnixMicro()
	b4 := a.UnixNano()
	t.Log(b, b2, b3, b4)

	// 秒数转time.Time
	bstr := time.Unix(b, 0)
	t.Log(bstr)

	//Time格式化为字符串
	c := a.Format("2006-01-02 15:04:05")
	t.Log(c)

	//随机睡眠时长
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))
}

func TestCompare(t *testing.T) {
	t1 := time.Now()
	t.Log(t1)

	t2 := time.Now()
	t.Log(t2)

	flag := t2.Equal(t1)
	t.Log("是否同一时间", flag)

	flag = t2.Before(t1)
	t.Log("t2 在 t1 之前", flag)

	flag = t2.After(t1)
	t.Log("t2 在 t1 之后", flag)
}

func Test1(t *testing.T) {
	str := "2022-04-06T08:09:03-04:00"
	i := strings.LastIndex(str, "-")
	t.Log(i)

	before := str[:i]
	after := str[i+1:]
	t.Log(before, after)

	beforeT, err := time.Parse("2006-01-02T15:04:05", before)
	t.Log(beforeT, err)

	afterD, err := time.ParseDuration(after)
	t.Log(afterD, err)
}

func Test2(t *testing.T) {
	now := time.Now()
	y := now.Year()
	d := now.Day()
	yd := now.YearDay()
	t.Log(y, d, yd)
}
