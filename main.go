package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
)

var ()

func main() {
	//macAdd, _ := exec.Command("./bin/bin").Output()
	macAdd, _ := exec.Command("./bin/invalid").Output()

	url := "http://localhost:9999"
	method := "GET"
	client := &http.Client{}
	toSend := strconv.Quote(string(macAdd))
	reqBody := map[string]any{
		"x-check": toSend,
	}

	fmt.Println(toSend)
	reqJSON, err := json.Marshal(reqBody)

	req, err := http.NewRequest(method, url, bytes.NewReader(reqJSON))

	if err != nil {
		fmt.Println(err)
		return
	}
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
