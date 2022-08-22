package json

import (
	"encoding/json"
	"testing"
)

type Media struct {
	Type string
	Url  string
}

func TestJson(t *testing.T) {
	media := Media{
		Type: "image",
		Url:  "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png",
	}
	jsonBytes, err := json.Marshal(media)
	t.Log(string(jsonBytes), err)

	var obj Media
	err = json.Unmarshal(jsonBytes, &obj)
	t.Log(obj, err)
}

func TestJsonNil(t *testing.T) {
	var ms1 []*Media
	ms2 := []*Media{}

	jsonBytes, err := json.Marshal(ms1)

	flag11 := ms1 == nil
	flag12 := string(jsonBytes) == ""
	t.Log(err, flag11, flag12, jsonBytes, string(jsonBytes))

	jsonBytes, err = json.Marshal(ms2)
	flag21 := ms1 == nil
	flag22 := string(jsonBytes) == ""
	t.Log(err, flag21, flag22, jsonBytes, string(jsonBytes))
}
