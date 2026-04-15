package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//go:embed input.json
var input []byte

type Input struct {
	Email []string `json:"email"`
}

func main() {
	//http://127.0.0.1:18011
	domain := os.Args[1]

	var obj Input
	err := json.Unmarshal(input, &obj)
	if err != nil {
		return
	}

	i := 0
	num := 500
	for i = 0; i < len(obj.Email)/num; i++ {
		tmp := Input{Email: obj.Email[i*num : (i+1)*num]}
		req(domain, tmp)
	}

	tmp := Input{Email: obj.Email[i*num : (i+1)*num]}
	req(domain, tmp)

}

func req(domain string, obj Input) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return
	}

	url := domain + "/ctrl/nothing-ear-stick/code/create"
	method := "POST"

	payload := bytes.NewReader(buf)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
