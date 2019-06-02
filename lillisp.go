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

type Pair struct {
  car interface{}
  cdr interface{}
}
type Atom string

func cons(a interface{}, b interface{}) *Pair {
  return &Pair{a, b}
}
func car(p *Pair) interface{} {
  return p.car
}
func cdr(p *Pair) interface{} {
  return p.cdr
}

func process_list(scanner *bufio.Scanner) *Pair {
  if (!scanner.Scan()) {
		panic("unmatched (") // we're in a list and ran out of tokens
  }
  token := scanner.Text()
  if (token == "(") {
    // start a new list and then add that to the one we're currently working on
    return cons(process_list(scanner), process_list(scanner))
  } else if (token == ")") {
    return nil
  } else {
    return cons(token, process_list(scanner))
  }
}

func process_item(scanner *bufio.Scanner) {
  if (!scanner.Scan()) {
    return // not sure what to do here, maybe panic?
  }
  token := scanner.Text()
  if (token == "(") {
    list := process_list(scanner)
    fmt.Printf("list: %v\n", list)
  } else if (token == ")") {
		panic("unmatched )")
  } else {
    fmt.Printf("atom: %v\n", token)
  }
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)
  process_item(scanner)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
}
