package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

// isSpace reports whether the character is a Unicode white space character.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

func scanWord(data []byte, start int, atEOF bool) (advance int, token []byte, err error) {
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[start+i:])
		if isSpace(r) || r == '(' || r == ')' {
			return start + i, data[start : start+i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	} else {
		return 0, nil, nil
	}
}

func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}

	// get a character
	r, width := utf8.DecodeRune(data[start:])

	switch r {
	case '(', ')':
		return start + width, data[start : start+width], nil
	default:
		return scanWord(data, start, atEOF)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
