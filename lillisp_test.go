package main

import "testing"

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
	input := "   a"
	width, token, err := scan([]byte(input), true)
	if width != 4 || string(token) != "a" || err != nil {
		t.Errorf("fail %v, %s, %v", width, token, err)
	}
}

func TestPrintItem(t *testing.T) {
	var p *Pair
	PrintItem(p) // should show ()

	p = Cons("bar", Nil)
	p = Cons("foo", p)
	PrintItem(p) // should show (foo bar)

	p = Cons(Cons("a", Nil), p)
	PrintItem(p) // should show ((a) foo bar)
}

func TestEq(t *testing.T) {
	p := Cons("qwe", Nil)

	tests := []struct {
		a, b, result interface{}
	}{
		{Nil, Nil, "t"},
		{"abc", "abc", "t"},
		{"abc", "ABC", Nil},
		{"abc", Nil, Nil},
		{Nil, "abc", Nil},
		{p, "abc", Nil},
		{p, p, Nil}, // they're the same but not an atom or nil
	}
	for _, test := range tests {
		result := Eq(test.a, test.b)
		if result != test.result {
			t.Errorf("fail Eq(%v, %v) == %v, not %v\n", test.a, test.b, result, test.result)
		}
	}
}

func TestAtom(t *testing.T) {
	p := Cons("qwe", Nil)

	tests := []struct {
		a, result interface{}
	}{
		{Nil, "t"},
		{"abc", "t"},
		{p, Nil},
	}
	for _, test := range tests {
		result := Atom(test.a)
		if result != test.result {
			t.Errorf("fail Atom(%v) == %v, not %v\n", test.a, result, test.result)
		}
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		a, result interface{}
	}{
		{Cons("quote", Cons("t", Nil)), "t"},
	}
	for _, test := range tests {
		result := Eval(test.a)
		if result != test.result {
			t.Errorf("fail Eval(%v) == %v, not %v\n", test.a, result, test.result)
		}
	}
}
