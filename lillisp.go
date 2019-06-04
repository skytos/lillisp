package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode/utf8"
)

// Atom ...
type Atom interface{}

// Pair represents the fundamental "cons-pair" in LISP-like languages
type Pair struct {
	car interface{}
	cdr interface{}
}

var funcMap = map[string]func(Atom, Atom) Atom{
	"+": func(a Atom, b Atom) Atom { return toInt(a) + toInt(b) },
	"-": func(a Atom, b Atom) Atom { return toInt(a) - toInt(b) },
	"*": func(a Atom, b Atom) Atom { return toInt(a) * toInt(b) },
	"/": func(a Atom, b Atom) Atom { return toInt(a) / toInt(b) },
}

func toInt(a Atom) int {
	if val, ok := a.(int); ok == true {
		return val
	} else if val, ok := a.(string); ok == true {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}

	log.Panicf("Could not convert %v to Number\n", a)
	return 0 // technically will never reach here
}

func toFunc(p interface{}) (func(a Atom, b Atom) Atom, bool) {
	op, ok := p.(string)
	if !ok {
		return nil, false
	}
	if funcMap[op] == nil {
		return nil, false
	}
	return funcMap[op], true
}

func toPair(p interface{}) (*Pair, bool) {
	pair, ok := p.(*Pair)
	if !ok {
		return nil, false
	}
	return pair, true
}

func eval(p interface{}) Atom {
	if p == nil {
		panic("Tried to eval nil!")
	}

	// p is an Atom, return it
	if _, success := toPair(p); success == false {
		return p
	}

	// if the first element is a function, recurse on the arguments
	if op, success := toFunc(Car(p)); success != false {
		return op(eval(Car(Cdr(p))), eval(Car(Cdr(Cdr(p)))))
	}

	panic("wtf are we doing here?!")
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
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) || r == '(' || r == ')' {
			return i, data[start:i], nil
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return 0, nil, nil
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

// Cons make a new Pair
func Cons(a interface{}, b interface{}) *Pair {
	return &Pair{a, b}
}

// Car get the first item of a Pair
func Car(p interface{}) interface{} {
	if pair, success := toPair(p); success != false {
		return pair.car
	}
	panic("Car called on non-pair")
}

// Cdr get the second item of a Pair
func Cdr(p interface{}) interface{} {
	if pair, success := toPair(p); success != false {
		return pair.cdr
	}
	panic("Cdr called on non-pair")
}

// processList make a list to represent the tokens being scanned
func processList(scanner *bufio.Scanner) *Pair {
	if !scanner.Scan() {
		panic("unmatched (") // we're in a list and ran out of tokens
	}
	token := scanner.Text()
	fmt.Printf("token: %s\n", token)
	if token == "(" {
		// start a new list and then add that to the one we're currently working on
		return Cons(processList(scanner), processList(scanner))
	} else if token == ")" {
		return nil
	} else {
		return Cons(token, processList(scanner))
	}
}

// processItem scan an item and print a representation of it, return the list structure if the item is a list,
// nil otherwise
func processItem(scanner *bufio.Scanner) *Pair {
	if !scanner.Scan() {
		// out of tokens
		return nil
	}
	token := scanner.Text()

	if token == "quit" {
		os.Exit(0)
	}

	if token == "(" {
		// start a new list
		list := processList(scanner)
		// and  print it
		fmt.Print("list: ")
		printItem(list)
		fmt.Println()
		return list
	} else if token == ")" {
		// we're not in a list so we shouldn't see ")"
		panic("unmatched )")
	} else {
		// must be an atom so print it
		fmt.Print("atom: ")
		printItem(token)
		fmt.Println()
		return nil
	}
}

// printItem print an item, duh
func printItem(i interface{}) {
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
	}
	if addSpace {
		fmt.Print(" ") // so there's a gap between items
	}
	printItem(Car(l))
	printList(Cdr(l).(*Pair), true) // print the rest, this assumes Cdr is a pair
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scan)

	// read print loop
	for {
		fmt.Print("lillisp> ")
		list := processItem(scanner)
		if list != nil {
			fmt.Println(eval(list))
		}
	}
}
