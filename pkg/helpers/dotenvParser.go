package helpers

import (
	"regexp"
	"strings"

	"github.com/0xF000D/goenv/pkg/utils"
)

var dotEnvLineRegex = regexp.MustCompile(
	`(?m)^\s*(?:export\s+)?([\w.-]+)(?:[^\S\r\n]*=[^\S\r\n]*|:\s+)` +
		`(\s*'(?:\\'|[^'])*'` + // single-quoted (multiline-safe)
		`|\s*"(?:\\"|[^"])*"` + // double-quoted (multiline-safe)
		`|\s*` + "`" + `(?:\\` + "`" + `|[^` + "`" + `])*` + "`" + // backtick-quoted (multiline-safe)
		`|[^#\r\n]*)` + // unquoted: any # starts a comment
		`[^\S\r\n]*(?:#[^\r\n]*)?[^\S\r\n]*$`, // trailing horizontal whitespace / inline comment
)

var windowsNewLine = regexp.MustCompile(`(?m)\r\n?`)
var newLine = regexp.MustCompile(`\\n`)
var carriageReturn = regexp.MustCompile(`\\r`)
var tab = regexp.MustCompile(`\\t`)

// Function to parse the any `.env` file
func DotEnvParser(src string, skipExpandForDoubleQuotes, skipConvertingWindowsNewlines, collectAllValues bool) map[string][]string {
	envObj := make(map[string][]string)

	if !skipConvertingWindowsNewlines {
		src = windowsNewLine.ReplaceAllString(src, "\n")
	}

	matches := dotEnvLineRegex.FindAllStringSubmatch(src, -1)

	for _, match := range matches {
		key := match[1]
		value := strings.TrimSpace(match[2])

		// Detect quoting style before any modification
		var maybeQuote byte
		if len(value) > 0 {
			maybeQuote = value[0]
		}

		// Remove surrounding quotes
		value = utils.Unquote(value)

		// Expand escape sequences only inside double-quoted values
		if maybeQuote == '"' && !skipExpandForDoubleQuotes {
			value = newLine.ReplaceAllString(value, "\n")
			value = carriageReturn.ReplaceAllString(value, "\r")
			value = tab.ReplaceAllString(value, "\t")
		}

		if collectAllValues {
			if _, ok := envObj[key]; ok {
				envObj[key] = append(envObj[key], value)
			} else {
				envObj[key] = []string{value}
			}
		} else {
			envObj[key] = []string{value}
		}
	}

	return envObj
}
