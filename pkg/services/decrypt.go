package services

import (
	"log"
	"os"
	"path/filepath"

	"github.com/0xF000D/goenv/pkg/helpers"
)

func Decrypt(filePaths []string, keysFilePath string) {
	if len(filePaths) == 1 && filePaths[0] == ".env" {
		_, err := os.Stat(".env")
		if err != nil {
			return
		}
	}

	// prased keys file
	var keysFileParsed map[string][]string

	if keysFilePath == ".env.keys" {
		_, err := os.Stat(".env.keys")
		if err == nil {
			// Read the keys file
			keysFileContent, err := os.ReadFile(keysFilePath)
			if err != nil {
				log.Fatalf("Unable to read file: %s", keysFilePath)
			}

			// parse the keys file
			keysFileParsed = helpers.DotEnvParser(string(keysFileContent), false, false, false)
		}
	} else {
		// Check if keys file exists
		helpers.CheckFileExists(keysFilePath)

		// Read the keys file
		keysFileContent, err := os.ReadFile(keysFilePath)
		if err != nil {
			log.Fatalf("Unable to read file: %s", keysFilePath)
		}

		// parse the keys file
		keysFileParsed = helpers.DotEnvParser(string(keysFileContent), false, false, false)
	}

	// loop over each file one by one
	for _, filePath := range filePaths {
		// Check if file exists
		helpers.CheckFileExists(filePath)

		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Unable to read file: %s", filePath)
		}

		// get the base name of file
		fileName := filepath.Base(filePath)

		// parse the file
		parsedEnvVars := helpers.DotEnvParser(string(content), false, false, false)

		// Check if this file contains any encrypted var
		containsEncryptedVar := helpers.ContainsEncryptedVar(parsedEnvVars)

		// Check if private exists for this file
		isPrivateKeyExistsForThisFile := helpers.IsPrivateKeyExistsForAFile(keysFileParsed, fileName)

		var privateKey string

		// if file contains encrypted var and private key does not exist
		if containsEncryptedVar && !isPrivateKeyExistsForThisFile {
			log.Fatalf("MISSING_PRIVATE_KEY: '%s' does not exist in file %s",
				helpers.PrivateKeyNameFromFileName(fileName), keysFilePath)
		} else if containsEncryptedVar && isPrivateKeyExistsForThisFile {
			// get the private key
			privateKey = helpers.GetPrivateKeyFromKeysFile(keysFileParsed, fileName)
		} else if !containsEncryptedVar && isPrivateKeyExistsForThisFile {
			return
		} else if !containsEncryptedVar && !isPrivateKeyExistsForThisFile {
			log.Fatalf("MISSING_PRIVATE_KEY: '%s' does not exist in file %s",
				helpers.PrivateKeyNameFromFileName(fileName), keysFilePath)
		}

		helpers.WriteEnvFileAfterDecryption(string(content), filePath, privateKey)
	}
}
