package main

import (
	"crypto/rsa"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"geocomply-lite/types"
	"geocomply-lite/utils"
)

var (
	pubKey = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA3r5CwPZdwC+ms5opTtny
+9PQ48334Gxc8Q+aMfOnUNLmOapClxxKFD5l8oJTjn+5UbZjKAKczkKRk1Baspl+
m694WlxYsJAcup6jir8dh9MmUWcAutH8Opy1phrPKl0KdvrcwAMCLGCamTUQF8x3
2k0Pbr0baBpyl55gELrZEo3ZEI1Ff0sMAfFPS4RwYIlttIutyRXcxFVnubG32nlo
wtIS6uoUvW3QgP264h1mmr5h8EFCWirVPFv41MnMmP72kSi2vo2UF7IDGfRXu0zQ
RiuszRjuC+GrnZx6Msffkzroopa2v9w06gcCCX+6XG/iwSZzvfhWBUStJXKYNnC+
J0WZoMX+DaIOmsUfWVqKrxubc9YvkJU0Z1Sv1KaFr4b7DbzJCpYRKLZ1NTADbjg5
ih5p4ieFHKgud+bomfRNyH+9tkfIYvVRfV5EWi+W6lnifTKdETFO7flMMvnEg8iz
+82tDLUvpBvNeaL1jwIst4gR8zeB56nc6IDx/tGAJhKTNjFbDxyrI9yAoOpndvCl
XC7v2VyKdf9Ufp+59EWgPvuzO/VfNxu6P8Hcfq+BJ7/Mqv56XRino0yTb3Xpa6tr
E0Pr76oDpJmx3/c+Q22vhinZ+fQg4J6ET4HDDBgym4ZhHCAiWGrCzGF/PfkKPc/b
db5PnM90AhxlLjwI/ImsLLsCAwEAAQ==
-----END RSA PUBLIC KEY-----`
	port             = ""
	version          = "1.0.0"
	getEncryptedPath = ""
	macAddresses     = []string{}
	parsedPubKey     *rsa.PublicKey
)

func init() {
	flag.StringVar(&port, "port", "8080", "allowed port for this agent")
	flag.StringVar(&getEncryptedPath, "get-encrypted-path", "/get-encrypted", "custom path to get encrypted data")
	flag.Parse()

	addresses, err := getMacAddr()
	if err != nil {
		log.Fatalf("cant get mac address %+v\n", err)
	}
	macAddresses = addresses
	parsedPubKey, err = utils.ParseRsaPublicKeyFromPemStr(pubKey)
	if err != nil {
		log.Fatal("cant parse public key")
	}

	if !strings.HasPrefix(getEncryptedPath, "/") {
		getEncryptedPath = "/" + getEncryptedPath
	}
}

func main() {
	http.HandleFunc(getEncryptedPath, func(w http.ResponseWriter, req *http.Request) {
		secret, err := toEncryptPayload()
		if err != nil {
			returnFailed(w)
			return
		}
		encrypted := utils.Encrypt(parsedPubKey, secret)
		out, err := buildOutput(string(encrypted))
		if err != nil {
			returnFailed(w)
			return
		}
		log.Println("Return encrypted data")
		w.Header().Add("Content-type", "application/json")
		w.Write(out)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Write([]byte(`{"status":"live"}`))
	})

	log.Println("Start http server on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("stopped server: %+v", err)
	}

}

func returnFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-type", "application/json")
	w.Write([]byte(`{"status":"failed"}`))

}

func toEncryptPayload() ([]byte, error) {
	return json.Marshal(types.EncryptPayload{
		Ts:           time.Now().UnixNano(),
		MacAddresses: macAddresses,
	})
}

func buildOutput(encrypted string) ([]byte, error) {
	quoted := strconv.Quote(encrypted)
	return json.Marshal(types.Response{
		Encrypted: quoted,
		Version:   version,
	})
}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var as []string
	for _, ifa := range ifas {
		addresses, err := ifa.Addrs()
		if err != nil {
			continue
		}
		if !containValidIP(addresses) {
			continue
		}
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func containValidIP(addresses []net.Addr) bool {
	for _, address := range addresses {
		// ignore ip loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return true
			}
		}
	}
	return false
}
