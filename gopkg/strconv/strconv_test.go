package strconv

import (
	"strconv"
	"testing"
)

// string 转 整型，浮点型
func TestConvInt(t *testing.T) {
	i := 1<<32 - 1
	//int转string 10/16进制用base表示(2<base<36)
	s1 := strconv.Itoa(i)
	s2 := strconv.FormatInt(int64(i), 10)
	s3 := strconv.FormatUint(uint64(i), 16)
	t.Log(i, s1, s2, s3)

	//string转int
	j1, err := strconv.Atoi(s1)
	if err != nil {
		t.Log(err)
	}
	//bitSize表示转换的位数，超出会报错。如2^8=256位，第一位为符号位
	j2, err := strconv.ParseInt(s3, 16, 33)
	if err != nil {
		t.Log(err)
	}

	t.Log(j1, j2)
}

func TestConvFloat(t *testing.T) {
	//float到string
	var f32 float32 = 3.1415926535
	t.Log(f32)
	/*
		fmt 可选格式如下
		'b' (-ddddp±ddd，二进制指数)
		'e' (-d.dddde±dd，十进制指数)
		'E' (-d.ddddE±dd，十进制指数)
		'f' (-ddd.dddd，没有指数)
		'g' ('e':大指数，'f':其它情况)
		'G' ('E':大指数，'f':其它情况)

		prec 控制精度（排除指数部分）
		fmt=f/e/E时，prec表示小数点后的数字个数
		fmt=g/G时，prec控制总的数字个数
		若prec=-1，表示使用最少数量的、但又必需的数字来表示f
	*/
	s1 := strconv.FormatFloat(float64(f32), 'E', -1, 32)
	s2 := strconv.FormatFloat(float64(f32), 'E', -1, 64)
	t.Log(s1, s2)

	//string到float32 精度会丢失
	f64, err := strconv.ParseFloat(s1, 32)
	t.Log(f64, err)
	//string到float64 精度不丢失
	f64, err = strconv.ParseFloat(s1, 64)
	t.Log(f64, err)
}
