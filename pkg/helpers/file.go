package helpers

import (
	"log"
	"os"

	"github.com/0xF000D/goenv/pkg/utils"
)

func CheckFileExists(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s: File does not exist: %s", utils.FILE_NOT_EXISTS, filePath)
		}
		log.Fatalf("%s: Unable to access file: %s", utils.FILE_NOT_ACCESSIBLE, filePath)
	}
}
