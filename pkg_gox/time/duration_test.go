package time

import (
	"testing"
	"time"
)

var TargetTime = getTargetTime()

func getTargetTime() time.Time {
	t, err := time.Parse("2006-01-02", "2050-01-01")
	if err != nil {
		panic(err)
	}
	return t
}

func TestSub(t *testing.T) {
	createTime := time.Now()
	sub := TargetTime.Sub(createTime)
	t.Log(sub, int64(sub))

	a := TargetTime.Unix() - createTime.Unix()
	t.Log(TargetTime.Unix(), createTime.Unix(), a)
}
