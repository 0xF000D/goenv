package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/0xF000D/goenv/pkg/utils"
	eciesgo "github.com/ecies/go/v2"
)

// function to encrypt value with ECIES
func EncryptValue(value string, publicKey *eciesgo.PublicKey) string {
	cipherText, err := eciesgo.Encrypt(publicKey, []byte(value))
	if err != nil {
		log.Fatalf("%s: Failed to encrypt value: %s", utils.ENCRYPTION_FAILED, err)
	}
	base64Encoded := base64.StdEncoding.EncodeToString(cipherText)
	return fmt.Sprintf("%s%s", utils.ENCRYPTED_VALUE_PREFIX, base64Encoded)
}

// function to encrypt value with AES256
func EncryptValueWithAES256(value string, key string) string {
	// Create AES block cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal("Failed to create AES Block cipher")
	}

	// Wrap AES block cipher in GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("Failed to wrap AES block cipher into GCM")
	}

	// Create a unique nonce
	// NOTE: never reuse nonce with the same key
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	// Encrypt the value
	cipherText := gcm.Seal(nonce, nonce, []byte(value), nil)

	return fmt.Sprintf("%s%s", utils.ENCRYPTED_VALUE_PREFIX, string(cipherText))
}
