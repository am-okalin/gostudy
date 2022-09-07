// Server2 is a minimal "echo" and counter server.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/webhook/orders", orders)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(r.URL, string(buf))
}

func orders(w http.ResponseWriter, r *http.Request) {
	info := strings.Builder{}

	info.WriteString(r.URL.String())
	info.WriteString("--------------------\n")

	for key, value := range r.Header {
		info.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	info.WriteString(string(buf))
	info.WriteString("\n")

	err = writeLog(info.String())
	if err != nil {
		log.Println(err)
	}
}

func writeLog(str string) error {
	const f = "./cmd./bizp/bizp.log"

	fh, err := os.OpenFile(f, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	defer fh.Close()

	_, err = fh.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}
