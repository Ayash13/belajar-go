package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenerateReferenceNo() string {
	now := time.Now()
	return fmt.Sprintf("%s%010d", now.Format("20060102"), rand.Int63n(10000000000))
}

func GenerateAuthCode() string {
	b := make([]byte, 64)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateAPIKey() string {
	return fmt.Sprintf("%s-%s-%s",
		randomHex(2),
		randomHex(2),
		randomHex(5),
	)
}

func GenerateAccountID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		randomHex(4),
		randomHex(2),
		randomHex(2),
		randomHex(2),
		randomHex(6),
	)
}

func VerifySignature(secretKey, payload, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func GenerateSignature(secretKey, payload string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func randomHex(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = "0123456789ABCDEF"[rand.Intn(16)]
	}
	return string(b)
}
