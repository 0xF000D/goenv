package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/0xF000D/goenv/pkg/utils"
)

// This function checks if the private for the given file
// exists in keys file?
func IsPrivateKeyExistsForAFile(parsedKeys map[string][]string, fileName string) bool {
	privateKeyName := PrivateKeyNameFromFileName(fileName)
	if _, ok := parsedKeys[privateKeyName]; ok {
		return true
	}
	return false
}

func GetPrivateKeyFromKeysFile(parsedKeys map[string][]string, fileName string) string {
	privateKeyName := PrivateKeyNameFromFileName(fileName)
	return parsedKeys[privateKeyName][0]
}

// Write a new private key to keys file
func WritePrivateKeyToKeysFile(privateKey string, envFileName string, keysFilePath string) {
	file, err := os.OpenFile(keysFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%s: Failed to open keys file: %s", utils.OPEN_FILE_FAILED, keysFilePath)
	}
	defer file.Close()

	// generate private key name from env file name
	privateKeyName := PrivateKeyNameFromFileName(envFileName)

	// new key to add
	newKey := fmt.Sprintf("\n# %s\n%s=%s\n", envFileName, privateKeyName, privateKey)

	// Create write
	writer := bufio.NewWriter(file)

	// Write the new key
	writer.WriteString(newKey)

	// flush the content to the file
	err = writer.Flush()
	if err != nil {
		log.Fatalf("%s: failed to write to the file: %s", utils.WRITE_FILE_FAILED, keysFilePath)
	}
}
