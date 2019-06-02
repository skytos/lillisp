package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

type Pair struct {
	car interface{}
	cdr interface{}
}

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

// make a new Pair
func cons(a interface{}, b interface{}) *Pair {
	return &Pair{a, b}
}

// get the first item of a Pair
func car(p *Pair) interface{} {
	return p.car
}

// get the second item of a Pair
func cdr(p *Pair) interface{} {
	return p.cdr
}

// make a list to represent the tokens being scanned
func process_list(scanner *bufio.Scanner) *Pair {
	if !scanner.Scan() {
		panic("unmatched (") // we're in a list and ran out of tokens
	}
	token := scanner.Text()
	if token == "(" {
		// start a new list and then add that to the one we're currently working on
		return cons(process_list(scanner), process_list(scanner))
	} else if token == ")" {
		return nil
	} else {
		return cons(token, process_list(scanner))
	}
}

// scan an item and print a representation of it, return true if more to do
func process_item(scanner *bufio.Scanner) bool {
	if !scanner.Scan() {
		// out of tokens
		return false
	}
	token := scanner.Text()
	if token == "(" {
		// start a new list
		list := process_list(scanner)
		// and  print it
		fmt.Print("list: ")
		print_item(list)
		fmt.Println()
	} else if token == ")" {
		// we're not in a list so we shouldn't see ")"
		panic("unmatched )")
	} else {
		// must be an atom so print it
		fmt.Print("atom: ")
		print_item(token)
		fmt.Println()
	}
	return true
}

// print an item, duh
func print_item(i interface{}) {
	p, is_pair := i.(*Pair)
	if is_pair {
		fmt.Print("(")
		print_list(p, false)
		fmt.Print(")")
	} else {
		// must be an atom
		a := i.(string)
		fmt.Printf("%v", a)
	}
}

// prints the contents of a list
// unless you're recursing pass false for add_space
func print_list(l *Pair, add_space bool) {
	if l == nil {
		// we're at the end of the list
		return
	} else {
		if add_space {
			fmt.Print(" ") // so there's a gap between items
		}
		print_item(car(l))
		print_list(cdr(l).(*Pair), true) // print the rest, this assumes cdr is a pair
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)
	// read print loop
	for process_item(scanner) {
	}
}
