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
	var ms1 []*Media
	ms2 := []*Media{}

	marshal, err := json.Marshal(ms1)
	flag11 := ms1 == nil
	flag12 := string(marshal) == ""
	t.Log(err, flag11, flag12, marshal, string(marshal))

	marshal, err = json.Marshal(ms2)
	flag21 := ms1 == nil
	flag22 := string(marshal) == ""
	t.Log(err, flag21, flag22, marshal, string(marshal))
}
