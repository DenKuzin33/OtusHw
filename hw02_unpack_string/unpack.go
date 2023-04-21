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

	for _, v := range runes {
		if unicode.IsDigit(v) {
			if currentChar == "" {
				return "", ErrInvalidString
			}

			charCount, _ := strconv.Atoi(string(v))
			sb.WriteString(strings.Repeat(currentChar, charCount))
			currentChar = ""
		} else {
			if currentChar != "" {
				sb.WriteString(currentChar)
			}
			currentChar = string(v)
		}
	}

	if !unicode.IsDigit(runes[len(runes)-1]) {
		sb.WriteString(string(runes[len(runes)-1]))
	}

	return sb.String(), nil
}
