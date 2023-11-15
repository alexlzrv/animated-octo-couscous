package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "SERVER PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile("./private.pem", privateKeyPEM, 0644)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "AGENT PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile("./public.pem", publicKeyPEM, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
