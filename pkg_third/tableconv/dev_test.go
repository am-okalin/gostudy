package dev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/am-okalin/kit/tableconv"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test1(t *testing.T) {
	fde, err := tableconv.Csv2Table(de, ',')
	if err != nil {
		t.Error(err)
	}
	fes, err := tableconv.Csv2Table(es, ',')
	if err != nil {
		t.Error(err)
	}
	feu, err := tableconv.Csv2Table(eu, ',')
	if err != nil {
		t.Error(err)
	}
	ffr, err := tableconv.Csv2Table(fr, ',')
	if err != nil {
		t.Error(err)
	}

	// 去重
	m := tableconv.Unique("imei1", fde, fes, feu, ffr)

	// 导出csv
	table := make([][]string, 0, len(m))
	for _, row := range m {
		table = append(table, row)
	}

	err = tableconv.ToCsv(table, all)
	if err != nil {
		t.Error(err)
	}
}

func Test2(t *testing.T) {
	fde, err := tableconv.Csv2Table(de, ',')
	if err != nil {
		t.Error(err)
	}
	fes, err := tableconv.Csv2Table(es, ',')
	if err != nil {
		t.Error(err)
	}
	feu, err := tableconv.Csv2Table(eu, ',')
	if err != nil {
		t.Error(err)
	}
	ffr, err := tableconv.Csv2Table(fr, ',')
	if err != nil {
		t.Error(err)
	}

	// 去重
	m := tableconv.Unique("imei1", fde, fes, feu, ffr)

	// 导出csv
	tableAll := make([][]string, 0, len(m))
	for _, row := range m {
		tableAll = append(tableAll, row)
	}

	// 切割CSV
	tables := TableCut(tableAll, 1000)

	//请求接口
	for _, table := range tables {
		err = mc(table)
		if err != nil {
			t.Error(err)
		}
	}
}

func TableCut(table [][]string, num int) [][][]string {
	tables := make([][][]string, 0)

	l := len(table)
	var i int
	for i = 0; i < l/num; i++ {
		tables = append(tables, table[i*num:(i+1)*num])
	}
	tables = append(tables, table[i*num:])

	return tables
}

type Equipment struct {
	Imei1          string
	ActivationTime string
}

func mc(table [][]string) error {
	objs := make([]Equipment, len(table))
	for i, row := range table {
		objs[i].Imei1 = row[0]
		objs[i].ActivationTime = row[1]
	}

	buf, err := json.Marshal(objs)
	if err != nil {
		return err
	}

	url := "https://directus.bbxlk.cc/items/fesfp_equipment_verify"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
