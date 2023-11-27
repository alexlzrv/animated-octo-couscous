package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
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

func DecryptMiddleware(key *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, key, body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(plaintext))
			next.ServeHTTP(w, r)
		})
	}
}
