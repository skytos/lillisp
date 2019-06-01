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

func scanWord(data []byte) (advance int, token []byte, err error) {
	// Scan until space, marking end of word.
	// We know that the next character is not a space!
	for width, i := 0, 0; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return i, data[0:i], nil
		} else if r == '(' || r == ')' {
			return i, data[0:i], nil
		}
	}
	// THIS IS PROBABLY WRONG
	return len(data), data[0:], nil
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
	if start > 0 {
		return start, nil, nil
	}

	// get a character
	r, width := utf8.DecodeRune(data)

	switch r {
	case '(', ')':
		return width, data[0:width], nil
	default:
		return scanWord(data)
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func main() {
	fmt.Printf("lisp 123\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
