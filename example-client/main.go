package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"geocomply-lite/types"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("http://localhost:8080/get-encrypted")

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var payload types.Response
	json.Unmarshal(body, &payload)

	url := "http://localhost:9999"
	method := "GET"
	client := &http.Client{}

	reqBody := map[string]any{
		"x-check": payload.Encrypted,
	}
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

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))
}
