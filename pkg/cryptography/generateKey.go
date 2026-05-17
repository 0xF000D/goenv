package cryptography

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"

	"github.com/0xF000D/goenv/pkg/utils"
	eciesgo "github.com/ecies/go/v2"
)

// this function generated random secret key for AES-256
func GenerateSecretKeyForAES() string {
	// 32 bytes for AES-256
	// 16 bytes for AES-128 and 24 bytes for AES-192
	key := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Fatal("Failed to generate Secret key for AES")
	}

	return hex.EncodeToString(key)
}

// function to generate key pair for ECIES algorithm
func GenerateKeyForECIES() *eciesgo.PrivateKey {
	k, err := eciesgo.GenerateKey()
	if err != nil {
		log.Fatalf("%s: Failed to generate key for ECIES", utils.KEY_GENERATION_FAILED)
	}
	return k
}

// Generate private key from existing private key
func GenerateKeyForECIESFromExistingKey(privateKey string) *eciesgo.PrivateKey {
	k, err := eciesgo.NewPrivateKeyFromHex(privateKey)
	if err != nil {
		log.Fatalf("%s: Failed to generate ECIES private key from existing private key!", utils.KEY_GENERATION_FAILED)
	}
	return k
}
