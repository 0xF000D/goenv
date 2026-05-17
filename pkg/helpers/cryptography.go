package helpers

import (
	"regexp"
	"strings"

	"github.com/0xF000D/goenv/pkg/utils"
)

var encryptionPattern = regexp.MustCompile(`^encrypted:.+`)

// check if the value is encrypted
func IsEncrypted(value string) bool {
	return encryptionPattern.MatchString(value)
}

func RemoveEncryptionPrefix(value string) string {
	value, _ = strings.CutPrefix(value, utils.ENCRYPTED_VALUE_PREFIX)
	return value
}
