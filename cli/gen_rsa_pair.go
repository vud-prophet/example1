package main

import (
	"geocomply-lite/utils"
	"os"
)

func main() {
	priv, pub := utils.GenerateRsaKeyPair()

	privateKeyFile, err := os.OpenFile("./private_key", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = utils.ExportRsaPrivateKeyAsPemStrToFile(privateKeyFile, priv)
	if err != nil {
		panic(err)
	}

	publicKeyFile, err := os.OpenFile("./public_key", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = utils.ExportRsaPublicKeyAsPemStrToFile(publicKeyFile, pub)
	if err != nil {
		panic(err)
	}
}
