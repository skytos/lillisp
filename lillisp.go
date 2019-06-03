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

var Nil *Pair

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
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) || r == '(' || r == ')' {
			return i, data[start:i], nil
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
func Cons(a interface{}, b interface{}) *Pair {
	return &Pair{a, b}
}

// get the first item of a Pair
func Car(p *Pair) interface{} {
	return p.car
}

// get the second item of a Pair
func Cdr(p *Pair) interface{} {
	return p.cdr
}

// returns atom "t" when a and b are both Nil or the same atom otherwise returns Nil
func Eq(a, b interface{}) interface{} {
	ap, aIsPair := a.(*Pair)
	bp, bIsPair := b.(*Pair)
	if aIsPair && bIsPair && ap == nil && bp == nil { // both Nil
		return "t"
	} else if !aIsPair && !bIsPair && a.(string) == b.(string) { // same atom
		return "t"
	} else {
		return Nil
	}
}

// returns "t" when a is an atom or Nil
func Atom(a interface{}) interface{} {
	ap, isPair := a.(*Pair)
	if isPair && ap == Nil || !isPair {
		return "t"
	} else {
		return Nil
	}
}

// evaluates exp
func Eval(exp interface{}) interface{} {
	return "t"
}

// expects a list of options where each option is a list of 2 items:
// the first being a condition and the second the expression to evaluate
// if the condition evalutes to true
func Cond(options *Pair) interface{} {
	if options == Nil {
		panic("Cond: no true conditions")
	}
	test, isPair := Eval(Car(Car(options).(*Pair))).(*Pair)
	if !isPair || test != Nil {
		return Eval(Car(Cdr(Car(options).(*Pair)).(*Pair)))
	} else {
		return Cond(Cdr(options).(*Pair))
	}
}

// make a list to represent the tokens being scanned
func processList(scanner *bufio.Scanner) *Pair {
	if !scanner.Scan() {
		panic("unmatched (") // we're in a list and ran out of tokens
	}
	token := scanner.Text()
	if token == "(" {
		// start a new list and then add that to the one we're currently working on
		return Cons(processList(scanner), processList(scanner))
	} else if token == ")" {
		return nil
	} else {
		return Cons(token, processList(scanner))
	}
}

// scan an item and print a representation of it, return true if more to do
func processItem(scanner *bufio.Scanner) bool {
	if !scanner.Scan() {
		// out of tokens
		return false
	}
	token := scanner.Text()
	if token == "(" {
		// start a new list
		list := processList(scanner)
		// and  print it
		fmt.Print("list: ")
		PrintItem(list)
		fmt.Println()
	} else if token == ")" {
		// we're not in a list so we shouldn't see ")"
		panic("unmatched )")
	} else {
		// must be an atom so print it
		fmt.Print("atom: ")
		PrintItem(token)
		fmt.Println()
	}
	return true
}

// print an item, duh
func PrintItem(i interface{}) {
	p, isPair := i.(*Pair)
	if isPair {
		fmt.Print("(")
		printList(p, false)
		fmt.Print(")")
	} else {
		// must be an atom
		a := i.(string)
		fmt.Printf("%v", a)
	}
}

// prints the contents of a list
// unless you're recursing pass false for addSpace
func printList(l *Pair, addSpace bool) {
	if l == nil {
		// we're at the end of the list
		return
	} else {
		if addSpace {
			fmt.Print(" ") // so there's a gap between items
		}
		PrintItem(Car(l))
		printList(Cdr(l).(*Pair), true) // print the rest, this assumes Cdr is a pair
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)
	// read print loop
	for processItem(scanner) {
	}
}
