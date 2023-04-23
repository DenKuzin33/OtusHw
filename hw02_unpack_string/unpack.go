package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(source string) (string, error) {
	var sb strings.Builder
	var currentChar string

	runes := []rune(source)

	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	escaped := false
	for _, v := range runes {
		switch {
		case unicode.IsDigit(v):
			if escaped {
				currentChar = string(v)
				escaped = false
				continue
			}

			if currentChar == "" {
				return "", ErrInvalidString
			}

			charCount, _ := strconv.Atoi(string(v))
			sb.WriteString(strings.Repeat(currentChar, charCount))
			currentChar = ""
		case v == 92:
			if currentChar != "" {
				sb.WriteString(currentChar)
				currentChar = ""
			}
			if escaped {
				if currentChar == `\` {
					sb.WriteString(`\`)
				} else {
					currentChar = `\`
				}
				escaped = false
			} else {
				escaped = true
			}
		default:
			if escaped {
				return "", ErrInvalidString
			}
			if currentChar != "" {
				sb.WriteString(currentChar)
			}
			currentChar = string(v)
		}
	}

	if currentChar != "" {
		sb.WriteString(currentChar)
	}

	return sb.String(), nil
}
