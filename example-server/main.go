package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"geocomply-lite/types"
	"geocomply-lite/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var encodedPrivate = ""

func init() {
	encodedPrivate = os.Getenv("PRIVATE_KEY")
}

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

	priv, err := base64.StdEncoding.DecodeString(encodedPrivate)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tDecryt, _ := strconv.Unquote(out)
	parasedPrivKey, err := utils.ParseRsaPrivateKeyFromPemStr(string(priv))
	if err != nil {
		fmt.Println(err)
		return false
	}

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
