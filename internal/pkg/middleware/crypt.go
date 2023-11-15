package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io"
	"net/http"
	"os"
)

func CryptMiddleware(signKey []byte) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if signKey == nil {
				return
			}

			if r.Header.Get("HashSHA256") == "" {
				next.ServeHTTP(w, r)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if len(body) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			h := hmac.New(sha256.New, signKey)
			h.Write(body)
			serverHash := h.Sum(nil)

			hash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if !hmac.Equal(hash, serverHash) {
				http.Error(w, "Invalid HashSHA256 header value", http.StatusBadRequest)
				return
			}

			w.Header().Set("HashSHA256", hex.EncodeToString(serverHash))

			next.ServeHTTP(w, r)
		})
	}
}

func DecryptMiddleware(privateKeyPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privateKeyPath == "" {
				return
			}

			privateKeyPEM, err := os.ReadFile(privateKeyPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			privateKeyBlock, _ := pem.Decode(privateKeyPEM)
			privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(plaintext))
			next.ServeHTTP(w, r)
		})
	}
}
