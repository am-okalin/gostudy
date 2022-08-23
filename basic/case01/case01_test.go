package case01

import (
	"log"
	"testing"
)

//type HumanParams struct {
//	Tall   float64
//	Weight float64
//}

type Human struct {
	//name   string
	//age    int
	//Params HumanParams
}

type case01 interface {
	//HumanSlice 获取人类的json字符串
	HumanSlice(str string) (list []Human, err error)
	//HumanWeightMap 获取人类的名称与体重的映射
	HumanWeightMap(list []Human) (wm map[string]float64, err error)
}

func testCase01(c case01) error {
	str := `[{"name":"downing","age":10,"params":{"tall":177.5,"weight":82.5}},{"name":"gary","age":10,"params":{"tall":178.5,"weight":83.5}}]`
	list, err := c.HumanSlice(str)
	if err != nil {
		return err
	}
	wm, err := c.HumanWeightMap(list)
	if err != nil {
		return err
	}
	log.Println(wm)
	return nil
}

func TestDowningCase01(t *testing.T) {
	var c DowningCase01
	err := testCase01(c)
	if err != nil {
		t.Log(err)
	}
}

type DowningCase01 struct {
}

func (d DowningCase01) HumanSlice(str string) (list []Human, err error) {
	panic("implement me")

	//err = json.Unmarshal([]byte(str), &list)
	//if err != nil {
	//	return nil, err
	//}
	//return list, nil
}

func (d DowningCase01) HumanWeightMap(list []Human) (wm map[string]float64, err error) {
	panic("implement me")

	//wm = make(map[string]float64)
	//
	//for _, row := range list {
	//	wm[row.name] = row.Params.Weight
	//}
	//
	//return wm, nil
}
