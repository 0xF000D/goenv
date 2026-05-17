package cryptography

import (
	"encoding/base64"
	"log"

	"github.com/0xF000D/goenv/pkg/utils"
	eciesgo "github.com/ecies/go/v2"
)

func DecryptValue(value string, privateKey string) string {
	k, err := eciesgo.NewPrivateKeyFromHex(privateKey)
	if err != nil {
		log.Fatalf("%s: The private key may be invalid!", utils.INVALID_PRIVATE_KEY)
	}

	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Fatalf("%s: Failed to decode value", utils.BASE64_DECODE_FAILED)
	}

	plainText, err := eciesgo.Decrypt(k, []byte(decodedValue))
	if err != nil {
		log.Fatalf("%s: Failed to decrypt the value, %s", utils.DECRYPTION_FAILED, err)
	}

	return string(plainText)
}

func DecryptValue2(value string, privateKey *eciesgo.PrivateKey) string {
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Fatalf("Failed to decode value")
	}
	plainText, err := eciesgo.Decrypt(privateKey, []byte(decodedValue))
	if err != nil {
		log.Fatalf("%s: Failed to decrypt the value, %s", utils.DECRYPTION_FAILED, err)
	}

	return string(plainText)
}
