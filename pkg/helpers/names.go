package helpers

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/0xF000D/goenv/pkg/utils"
)

// This function returns the name of file to which
// given private key belongs
func FileNameFromPrivateKeyName(privateKeyName string) (string, error) {
	fileNameHexEncoded := privateKeyName[len(fmt.Sprintf("%s_", utils.PRIVATE_KEY_NAME_PREFIX)):]
	decoded, err := hex.DecodeString(fileNameHexEncoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// This functions generates the private key name for the
// given env file name
func PrivateKeyNameFromFileName(fileName string) string {
	return fmt.Sprintf("%s_%s", utils.PRIVATE_KEY_NAME_PREFIX, strings.ToUpper(hex.EncodeToString([]byte(fileName))))
}

// This functions generates the public key name for the
// given env file name
func PublicKeyNameFromFileName(fileName string) string {
	return fmt.Sprintf("%s_%s", utils.PUBLIC_KEY_NAME_PREFIX, strings.ToUpper(hex.EncodeToString([]byte(fileName))))
}
