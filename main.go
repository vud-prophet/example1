package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	priv = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	pub = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAv/AJHsoL/LHtwehGPEAd
i41m79hHY7+pOl7BH2qq/ZZNkCXWJ9M6/Xdd4igC1L+BtVZluEiCN/lgV+QIvDy9
+mFMcmB9LOMzOmzCvqS5VNF9w2eLyE0SdwPu20v7mN1+VRx8X0BdFJYksbmkKI+9
wz81A3DiqNo1yYDAInBp2T6zzVbyfYG786wSF5JkqdKvib6jCmjB+Y9NABKaIBLu
ITKD6jVbHGueI3QxeihJwL3ifFtmeYsXfMBnNsvCVVHOiu9UgD6pFHf0h+FfFeTi
959T9xNmOTd5DlAFOVA8PJ7PmuNkazXMv1l5Nst7aN9jnPLdGdjLkY3747WagJ/u
xFyyS7BxPcy19lC31JSA01rEnmko8Qdwp7pOWLWcYCNt4FbaN3R87kKgbcGdn4Z7
q1D5TqJQZ+Z8lqvSap4QEsd8KlyY3Zkd9yZzaVj3MNFap2fCMj9o/o1+GfL+uumu
yoAV9aox3rJcceGcFANN1jZiY4iX2EYiEFZmzBXUSbY3k9GAWwEFvbUVTejGZ/zC
eJZ+s7JhGylypj5pz53sN10GQ0Nza8sJl6deIBv8wK7vclDXgGhKHa4etKVIRPe/
eDXYpYusFSuFga1IUDGsyI9nqHYEGnfu4QzeQSxeg6Jd+SpUC/nTnNUqfWhvf9z6
JS53C1gGiZIVXK3ahPKQD/UCAwEAAQ==
-----END RSA PUBLIC KEY-----`
)

var whitelisted = map[string]struct{}{"40:23:43:c5:08:0b": {}}

func main() {
	out, _ := exec.Command("./bin/bin").Output()

	parasedPrivKey, _ := ParseRsaPrivateKeyFromPemStr(priv)
	decrypted := Decrypt(out, parasedPrivKey)
	outAddress := string(decrypted)
	fmt.Println("decrypted ", outAddress)

	parts := strings.Split(outAddress, "_")
	if len(parts) != 2 {
		fmt.Println("Malformed")
		return
	}
	addresses := strings.Split(parts[1], ",")

	for _, address := range addresses {
		if _, ok := whitelisted[address]; ok {
			fmt.Println("You shall pass")
			break
		}
	}
}
func Decrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) []byte {
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}
	return decryptedBytes
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	}
	return nil, errors.New("Key type is not RSA")
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
