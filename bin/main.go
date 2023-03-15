package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

var pubKey = `-----BEGIN RSA PUBLIC KEY-----
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

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func main() {
	as, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
	}

	//priv, pub := GenerateRsaKeyPair()

	//privPem := ExportRsaPrivateKeyAsPemStr(priv)
	//pubPem, _ := ExportRsaPublicKeyAsPemStr(pub)

	//fmt.Println(privPem)
	//fmt.Println(pubPem)

	parsedPubKey, _ := ParseRsaPublicKeyFromPemStr(pubKey)

	ec := Encrypt(parsedPubKey, strings.Join(as, ","))
	fmt.Print(string(ec))
}

func Encrypt(publicKey *rsa.PublicKey, addresses string) []byte {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(addresses),
		nil)
	if err != nil {
		panic(err)
	}
	return encryptedBytes
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
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
