package agent

import (
	"crypto/rand"
	"crypto/rsa"
)

func encrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
