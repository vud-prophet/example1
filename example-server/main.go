package main

import (
	"encoding/json"
	"fmt"
	"geocomply-lite/types"
	"geocomply-lite/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	priv = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if check(r) {
			fmt.Fprint(w, "OK")
			return
		}
		fmt.Fprint(w, "Stop")
		return
	})

	http.ListenAndServe(":9999", mux)
}

var whitelisted = map[string]struct{}{"02:42:9f:4c:ac:dd": {}}

func check(r *http.Request) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	reqBody := map[string]any{}
	json.Unmarshal(body, &reqBody)
	out := reqBody["x-check"].(string)

	tDecryt, _ := strconv.Unquote(out)
	parasedPrivKey, _ := utils.ParseRsaPrivateKeyFromPemStr(priv)
	decrypted, err := utils.Decrypt([]byte(tDecryt), parasedPrivKey)
	if err != nil {
		fmt.Println(err)
		return false
	}

	payload := types.EncryptPayload{}
	json.Unmarshal(decrypted, &payload)

	if len(payload.MacAddresses) == 0 {
		fmt.Println("Malformed")
		return false
	}
	for _, address := range payload.MacAddresses {
		if _, ok := whitelisted[address]; ok {
			fmt.Println("You shall pass")
			return true
		}
	}
	return false
}
