package helpers

import (
	"fmt"
	"os"
	"testing"
)

func TestDotEnvParser(t *testing.T) {
	// open the .env file
	content, err := os.ReadFile("../../files/.env")
	if err != nil {
		t.Error("Failed to read the .env file")
	}

	parsed := DotEnvParser(string(content), false, false, false)

	// test cases for each line
	tests := []struct {
		key      string
		expected string
		msg      string
	}{
		{"BASIC", "basic", "sets basic environment variable"},
		{"AFTER_LINE", "after_line", "reads after a skipped line"},
		{"EMPTY", "", "defaults empty values to empty string"},
		{"EMPTY_SINGLE_QUOTES", "", "defaults empty values to empty string"},
		{"EMPTY_DOUBLE_QUOTES", "", "defaults empty values to empty string"},
		{"EMPTY_BACKTICKS", "", "defaults empty values to empty string"},
		{"SINGLE_QUOTES", "single_quotes", "escapes single quoted values"},
		{"SINGLE_QUOTES_SPACED", "    single quotes    ", "respects surrounding spaces in single quotes"},
		{"DOUBLE_QUOTES", "double_quotes", "escapes double quoted values"},
		{"DOUBLE_QUOTES_SPACED", "    double quotes    ", "respects surrounding spaces in double quotes"},
		{"DOUBLE_QUOTES_INSIDE_SINGLE", "double \"quotes\" work inside single quotes", "respects double quotes inside single quotes"},
		{"DOUBLE_QUOTES_WITH_NO_SPACE_BRACKET", "{ port: $MONGOLAB_PORT}", "respects spacing for badly formed brackets"},
		{"SINGLE_QUOTES_INSIDE_DOUBLE", "single 'quotes' work inside double quotes", "respects single quotes inside double quotes"},
		{"BACKTICKS_INSIDE_SINGLE", "`backticks` work inside single quotes", "respects backticks inside single quotes"},
		{"BACKTICKS_INSIDE_DOUBLE", "`backticks` work inside double quotes", "respects backticks inside double quotes"},
		{"BACKTICKS", "backticks", ""},
		{"BACKTICKS_SPACED", "    backticks    ", ""},
		{"DOUBLE_QUOTES_INSIDE_BACKTICKS", "double \"quotes\" work inside backticks", "respects double quotes inside backticks"},
		{"SINGLE_QUOTES_INSIDE_BACKTICKS", "single 'quotes' work inside backticks", "respects single quotes inside backticks"},
		{"DOUBLE_AND_SINGLE_QUOTES_INSIDE_BACKTICKS", "double \"quotes\" and single 'quotes' work inside backticks", "respects single quotes inside backticks"},
		{"EXPAND_NEWLINES", "expand\nnew\nlines", "expands newlines but only if double quoted"},
		{"DONT_EXPAND_UNQUOTED", "dontexpand\\nnewlines", "expands newlines but only if double quoted"},
		{"DONT_EXPAND_SQUOTED", "dontexpand\\nnewlines", "expands newlines but only if double quoted"},
		{"INLINE_COMMENTS", "inline comments", "ignores inline comments"},
		{"INLINE_COMMENTS_SINGLE_QUOTES", "inline comments outside of #singlequotes", "respects # inside single quotes"},
		{"INLINE_COMMENTS_DOUBLE_QUOTES", "inline comments outside of #doublequotes", "respects # inside double quotes"},
		{"INLINE_COMMENTS_BACKTICKS", "inline comments outside of #backticks", "respects # inside backticks"},
		{"INLINE_COMMENTS_SPACE", "inline comments start with a", "treats # character as start of comment"},
		{"EQUAL_SIGNS", "equals==", "respects equals signs in values"},
		{"RETAIN_INNER_QUOTES", `{"foo": "bar"}`, "retains inner quotes"},
		{"RETAIN_INNER_QUOTES_AS_STRING", `{"foo": "bar"}`, "retains inner quotes"},
		{"RETAIN_INNER_QUOTES_AS_BACKTICKS", `{"foo": "bar's"}`, "retains inner quotes"},
		{"TRIM_SPACE_FROM_UNQUOTED", "some spaced out string", "retains spaces in string"},
		{"USERNAME", "therealnerdybeast@example.tld", "parses email addresses completely"},
		{"SPACED_KEY", "parsed", "parses keys and values surrounded by spaces"},
	}

	for _, tt := range tests {
		if val, ok := parsed[tt.key]; !ok || val[0] != tt.expected {
			t.Errorf("%s: expected %q, got %q", tt.key, tt.expected, val)
		}
	}

	if _, ok := parsed["COMMENTS"]; ok {
		t.Error("should ignore commented lines")
	}
}

func TestDotEnvParserMultiline(t *testing.T) {
	// open the .env file
	content, err := os.ReadFile("../../files/.env.multiline")
	if err != nil {
		t.Error("Failed to open the .env file")
	}

	parsed := DotEnvParser(string(content), false, false, false)

	// Multiline assertions
	tests := []struct {
		key      string
		expected string
		desc     string
	}{
		{
			"MULTI_DOUBLE_QUOTED",
			"THIS\nIS\nA\nMULTILINE\nSTRING",
			"parses multi line values in double quotes",
		},
		{
			"MULTI_SINGLE_QUOTED",
			"THIS\nIS\nA\nMULTILINE\nSTRING",
			"parses multi line values in single quotes",
		},
		{
			"MULTI_BACKTICKED",
			"THIS\nIS\nA\n\"MULTILINE'S\"\nSTRING",
			"parses multi line values in backticks",
		},
		{
			"MULTI_PEM_DOUBLE_QUOTED",
			"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnNl1tL3QjKp3DZWM0T3u\nLgGJQwu9WqyzHKZ6WIA5T+7zPjO1L8l3S8k8YzBrfH4mqWOD1GBI8Yjq2L1ac3Y/\nbTdfHN8CmQr2iDJC0C6zY8YV93oZB3x0zC/LPbRYpF8f6OqX1lZj5vo2zJZy4fI/\nkKcI5jHYc8VJq+KCuRZrvn+3V+KuL9tF9v8ZgjF2PZbU+LsCy5Yqg1M8f5Jp5f6V\nu4QuUoobAgMBAAE=\n-----END PUBLIC KEY-----",
			"parses multi line pem key in double quotes",
		},
	}

	for _, tt := range tests {
		fmt.Println(tt)
		if parsed[tt.key][0] != tt.expected {
			t.Errorf("%s: expected %q, got %q", tt.desc, tt.expected, parsed[tt.key][0])
		}
	}
}
