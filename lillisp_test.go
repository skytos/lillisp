package main

import (
	"bufio"
	"strings"
	"testing"
)

const checkMark = "\u2705"
const xMark = "\u274c"

func TestTokenizer(t *testing.T) {

	t.Log(`Testing basic tokenization`)
	{
		input := "()"
		var width int
		var token []byte
		var err error

		width, token, err = scan([]byte(input), false)
		if !(width == 1 && string(token) == "(" && err == nil) {
			t.Errorf("\t%v The first token was bad: %v, %s, %v", xMark, width, token, err)
		} else {
			t.Logf("\t%v  First token is as we expect!", checkMark)
		}

		width, token, err = scan([]byte(input[1:]), true)
		if !(width == 1 && string(token) == ")" && err == nil) {
			t.Errorf("\t%v The second token was bad: %v, %s, %v", xMark, width, token, err)
		} else {
			t.Logf("\t%v  Second token is as we expect!", checkMark)
		}

	}

	t.Log(`Testing moar complicated tokenization`)
	{
		input := []byte("(asd asd as)  (  ( asd a ))a)")
		var width int
		var token []byte
		var err error

		results := []struct {
			width int
			token string
		}{
			{1, "("},
			{3, "asd"},
			{4, "asd"},
			{3, "as"},
			{1, ")"},
			{3, "("},
			{3, "("},
			{4, "asd"},
			{2, "a"},
			{2, ")"},
			{1, ")"},
			{1, "a"},
			{1, ")"},
		}

		start := 0
		for i, testToken := range results {
			width, token, err = scan(input[start:], false)
			start += testToken.width

			if testToken.width != width || testToken.token != string(token) || err != nil {
				t.Errorf("\t%v(%v) Bad token returned! Width %v, Token %v, Error %v", xMark, i, width, string(token), err)
			}
		}
		t.Logf("\t%v  Pass!", checkMark)
	}
}

func TestTokenizerLeadingSpaces(t *testing.T) {
	t.Log(`Testing tokenizer nomming leading spaces`)
	{
		input := "   a"
		width, token, err := scan([]byte(input), true)
		if width != 4 || string(token) != "a" || err != nil {
			t.Errorf("\t%v fail %v, %s, %v", xMark, width, token, err)
		}
		t.Logf("\t%v  Pass!", checkMark)

	}
}

func TestAddition(t *testing.T) {
	t.Log(`Tetsing basic addition`)
	{
		input := "(+ 4 3)"
		testScanner := bufio.NewScanner(strings.NewReader(input))
		testScanner.Split(scan)
		result := eval(processItem(testScanner))
		if result != 7 {
			t.Errorf("\t%v Result not OK: %v", xMark, result)
		}
	}

	t.Log(`Testing nested addition`)
	{
		input := "(+ (+ 4 3) (+ 5 2))"
		testScanner := bufio.NewScanner(strings.NewReader(input))
		testScanner.Split(scan)
		result := eval(processItem(testScanner))
		if result != 14 {
			t.Errorf("\t%v Result not OK: %v", xMark, result)
		}
	}
}

func TestPrintItem(t *testing.T) {
	var p, Nil *Pair
	p = nil
	printItem(p) // should show ()

	p = Cons("bar", Nil)
	p = Cons("foo", p)
	printItem(p) // should show (foo bar)

	p = Cons(Cons("a", Nil), p)
	printItem(p) // should show ((a) foo bar) */
}
