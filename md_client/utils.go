package md_client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"
)

func generateRandomString(length int) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the random digits
	digits := make([]byte, length)

	// Generate random digits and add them to the slice
	for i := range digits {
		digits[i] = byte(rand.Intn(10) + '0')
	}

	// Convert the byte slice to a string and return it
	return string(digits)
}

func generateSignature(secret, sendingTime, messageType, messageSeqNumber, senderCompId, targetCompId, password string) string {
	hmacKey, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Println("Error decoding base64:", err)
		return ""
	}

	h := hmac.New(sha256.New, hmacKey)
	h.Write([]byte(strings.Join([]string{sendingTime, messageType, messageSeqNumber, senderCompId, targetCompId, password}, "\x01")))
	signature := h.Sum(nil)

	signB64 := base64.StdEncoding.EncodeToString(signature)
	return signB64
}
