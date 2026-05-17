package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xF000D/goenv/pkg/cryptography"
	"github.com/0xF000D/goenv/pkg/utils"
	eciesgo "github.com/ecies/go/v2"
)

// This function checks if the public for the given file
// exists in that file?
func IsPublicKeyExistsInEnvFile(envVars map[string][]string, fileName string) bool {
	publicKeyName := PublicKeyNameFromFileName(fileName)
	if _, ok := envVars[publicKeyName]; ok {
		return true
	}
	return false
}

// this function check if the given env file
// contains any encrypted env var
func ContainsEncryptedVar(envVars map[string][]string) bool {
	for _, values := range envVars {
		for _, value := range values {
			if strings.HasPrefix(value, utils.ENCRYPTED_VALUE_PREFIX) {
				return true
			}
		}
	}

	return false
}

// This function write the encrypted values of env vars to
// env file
func WriteEnvFile(content string, filepath string, publicKey *eciesgo.PublicKey) {
	// open the file
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("%s: Error opening file: %s", utils.OPEN_FILE_FAILED, filepath)
	}
	defer file.Close()

	// create write instance
	writer := bufio.NewWriter(file)

	// Find the matches
	content = windowsNewLine.ReplaceAllString(content, "\n")
	matches := dotEnvLineRegex.FindAllStringSubmatchIndex(content, -1)

	lastIndex := 0

	// handle each match
	for _, match := range matches {
		start := match[0]
		keyStart, keyEnd := match[2], match[3]
		valStart, valEnd := match[4], match[5]

		// if something is there after the value of this line ends and
		// before the start of the key of the next line
		if start > lastIndex {
			_, err := writer.WriteString(content[lastIndex:keyStart])
			if err != nil {
				fmt.Println("Error writing string:", err)
				return
			}
		}

		// Unquote the value
		value := utils.Unquote(content[valStart:valEnd])
		if !IsEncrypted(value) {
			if value == "" {
				value = utils.EMPTY_VALUE_PLACEHOLDER
			}
			value = cryptography.EncryptValue(value, publicKey)
		}

		// Key and value and if there is something between key end and equals sign
		// or something between equals sign and value start
		line := fmt.Sprintf("%s%s=%s%s", content[keyStart:keyEnd], strings.Split(content[keyEnd:valStart], "=")[0], strings.Split(content[keyEnd:valStart], "=")[1], value)
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Println("Error writing string:", err)
			return
		}

		lastIndex = valEnd
	}

	// if something is there where the last value ends and end of file
	if lastIndex < len(content) {
		_, err := writer.WriteString(content[lastIndex:])
		if err != nil {
			fmt.Println("Error writing string:", err)
			return
		}
	}

	// flush the content to the file
	err = writer.Flush()
	if err != nil {
		log.Fatalf("%s: failed to write to the file: %s", utils.WRITE_FILE_FAILED, filepath)
	}
}

func WriteEnvFileAfterDecryption(content string, filepath string, privateKey string) {
	var stringsToWrite []string

	// Find the matches
	content = windowsNewLine.ReplaceAllString(content, "\n")
	matches := dotEnvLineRegex.FindAllStringSubmatchIndex(content, -1)

	lastIndex := 0

	// handle each match
	for _, match := range matches {
		start := match[0]
		keyStart, keyEnd := match[2], match[3]
		valStart, valEnd := match[4], match[5]

		// if something is there after the value of this line ends and
		// before the start of the key of the next line
		if start > lastIndex {
			stringsToWrite = append(stringsToWrite, content[lastIndex:keyStart])
		}

		// Unquote the value
		value := utils.Unquote(content[valStart:valEnd])
		if IsEncrypted(value) {
			value = RemoveEncryptionPrefix(value)
			value = cryptography.DecryptValue(value, privateKey)
			if value == utils.EMPTY_VALUE_PLACEHOLDER {
				value = ""
			}
		}

		// Key and value and if there is something between key end and equals sign
		// or something between equals sign and value start
		line := fmt.Sprintf("%s%s=%s%s", content[keyStart:keyEnd], strings.Split(content[keyEnd:valStart], "=")[0], strings.Split(content[keyEnd:valStart], "=")[1], value)
		stringsToWrite = append(stringsToWrite, line)

		lastIndex = valEnd
	}

	// if something is there where the last value ends and end of file
	if lastIndex < len(content) {
		stringsToWrite = append(stringsToWrite, content[lastIndex:])
	}

	// open the file
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("%s: Error opening file: %s", utils.OPEN_FILE_FAILED, filepath)
	}
	defer file.Close()

	// create write instance
	writer := bufio.NewWriter(file)

	for _, stringToWrite := range stringsToWrite {
		_, err := writer.WriteString(stringToWrite)
		if err != nil {
			fmt.Println("Error writing string:", err)
			return
		}
	}

	// flush the content to the file
	err = writer.Flush()
	if err != nil {
		log.Fatalf("%s: failed to write to the file: %s", utils.WRITE_FILE_FAILED, filepath)
	}
}
