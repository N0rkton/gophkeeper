// Package utils provides utility functions
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"log"
)

var nonce = []byte{156, 123, 210, 167, 214, 230, 92, 233, 232, 233, 172, 192}

// Encrypt - use the AES cipher in Galois/Counter Mode (GCM) to perform authenticated encryption.
func Encrypt(text string, key []byte) string {
	// Generate a new AES cipher block using the secret key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes %v", err)
	}
	// Create a new GCM (Galois/Counter Mode) cipher using the AES block
	// GCM provides authenticated encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("GCM %v", err)
	}
	// Encrypt the plaintext using AES-GCM
	ciphertext := aesGCM.Seal(nil, nonce, []byte(text), nil)
	return base64.RawStdEncoding.EncodeToString(ciphertext)
}

// Decrypt - use the AES cipher in Galois/Counter Mode (GCM) to perform authenticated decryption.
func Decrypt(text string, key []byte) string {
	// Generate a new AES cipher block using the secret key
	decodedCiphertext, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		log.Fatalf("decrypt 1 %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new GCM (Galois/Counter Mode) cipher using the AES block
	// GCM provides authenticated encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("decrypt 2 %v", err)
	}
	// Decrypt the ciphertext using AES-GCM
	decrypted, err := aesGCM.Open(nil, nonce, decodedCiphertext, nil)
	if err != nil {
		log.Printf("decrypt 3 %v", err)
	}
	return string(decrypted)
}
