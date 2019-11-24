package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// parser with type checker

type tcType int

const (
	tyIllTyped tcType = 0
	tyInt      tcType = 1
	tyBool     tcType = 2
)

func showType(t tcType) string {
	var s string
	switch {
	case t == tyInt:
		s = "Int"
	case t == tyBool:
		s = "Bool"
	case t == tyIllTyped:
		s = "Illtyped"
	}
	return s
}

// AST

type exp interface {
	pretty() string
	infer() tcType
}

type err int
type num int
type boolean bool
type mult [2]exp
type plus [2]exp
type or [2]exp
type and [2]exp

// pretty print

func (e err) pretty() string {
	return "Syntax Error"
}

func (x num) pretty() string {
	return strconv.Itoa(int(x))
}

func (x boolean) pretty() string {
	return strconv.FormatBool(bool(x))
}

func (e mult) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "*"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e plus) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "+"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e or) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "||"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e and) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "&&"
	x += e[1].pretty()
	x += ")"

	return x
}

// Type inferencer/checker
func (e err) infer() tcType {
	return tyIllTyped
}

func (x num) infer() tcType {
	return tyInt
}

func (x boolean) infer() tcType {
	return tyBool
}

func (e mult) infer() tcType {
	t1 := e[0].infer()
	t2 := e[1].infer()
	if t1 == tyInt && t2 == tyInt {
		return tyInt
	}
	return tyIllTyped
}

func (e plus) infer() tcType {
	t1 := e[0].infer()
	t2 := e[1].infer()
	if t1 == tyInt && t2 == tyInt {
		return tyInt
	}
	return tyIllTyped
}

func (e or) infer() tcType {
	t1 := e[0].infer()
	t2 := e[1].infer()
	if t1 == tyBool && t2 == tyBool {
		return tyBool
	}
	return tyIllTyped
}

func (e and) infer() tcType {
	t1 := e[0].infer()
	t2 := e[1].infer()
	if t1 == tyBool && t2 == tyBool {
		return tyBool
	}
	return tyIllTyped
}

// Simple scanner/lexer

// Tokens
const (
	ERR   = -1
	EOS   = 0 // End of string
	ZERO  = 1
	ONE   = 2
	TWO   = 3
	OPEN  = 4
	CLOSE = 5
	PLUS  = 6
	MULT  = 7
	AND   = 8
	OR    = 9
	TRUE  = 10
	FALSE = 11
)

func scan(s string) (string, int) {
	for {
		switch {
		case len(s) == 0:
			// print("EOS")
			return s, EOS
		case s[0] == '0':
			// print("ZERO ")
			return s[1:len(s)], ZERO
		case s[0] == '1':
			// print("ONE ")
			return s[1:len(s)], ONE
		case s[0] == '2':
			// print("TWO ")
			return s[1:len(s)], TWO
		case s[0] == '+':
			// print("PLUS ")
			return s[1:len(s)], PLUS
		case s[0] == '*':
			// print("MULT ")
			return s[1:len(s)], MULT
		case s[0] == '(':
			// print("OPEN ")
			return s[1:len(s)], OPEN
		case s[0] == ')':
			// print("CLOSE ")
			return s[1:len(s)], CLOSE
		case len(s) >= 2 && s[0:2] == "&&":
			// print("AND ")
			return s[2:len(s)], AND
		case len(s) >= 2 && s[0:2] == "||":
			// print("OR ")
			return s[2:len(s)], OR
		case len(s) >= 4 && s[0:4] == "true":
			// print("TRUE ")
			return s[4:len(s)], TRUE
		case len(s) >= 5 && s[0:5] == "false":
			// print("FALSE ")
			return s[5:len(s)], FALSE
		case isSpace(s[0]):
			s = s[1:len(s)]
		default:
			// print("ERR ")
			return s[1:len(s)], ERR
		}

	}
}

func isSpace(b byte) bool {
	r, _ := utf8.DecodeRune([]byte{b})
	return unicode.IsSpace(r)
}

type state struct {
	s   *string
	tok int
}

func next(s *state) {
	s2, tok := scan(*s.s)

	s.s = &s2
	s.tok = tok
}

// EP  ::= EO EP2
func parseEP(s *state) (bool, exp) {
	b, e := parseEO(s)
	if !b {
		return false, e
	}
	return parseEP2(s, e)
}

// EP2 ::= + EO EP2 |
func parseEP2(s *state, e exp) (bool, exp) {
	if s.tok == PLUS {
		next(s)
		b, f := parseEO(s)
		if !b {
			return false, e
		}
		t := (plus)([2]exp{e, f})
		return parseEP2(s, t)
	}

	return true, e
}

// EO  ::= EM EO2
func parseEO(s *state) (bool, exp) {
	b, e := parseEM(s)
	if !b {
		return false, e
	}
	return parseEO2(s, e)
}

// EO2 ::= || EM EO2 |
func parseEO2(s *state, e exp) (bool, exp) {
	if s.tok == OR {
		next(s)
		b, f := parseEM(s)
		if !b {
			return false, e
		}
		t := (or)([2]exp{e, f})
		return parseEO2(s, t)
	}
	return true, e
}

// EM  ::= EA EM2
func parseEM(s *state) (bool, exp) {
	b, e := parseEA(s)
	if !b {
		return false, e
	}
	return parseEM2(s, e)
}

// EM2 ::= * EA EM2 |
func parseEM2(s *state, e exp) (bool, exp) {
	if s.tok == MULT {
		next(s)
		b, f := parseEA(s)
		if !b {
			return false, e
		}
		t := (mult)([2]exp{e, f})
		return parseEM2(s, t)
	}
	return true, e
}

// EA  ::= F EA2
func parseEA(s *state) (bool, exp) {
	b, e := parseF(s)
	if !b {
		return false, e
	}
	return parseEA2(s, e)
}

// EA2 ::= && F EA2 |
func parseEA2(s *state, e exp) (bool, exp) {
	if s.tok == AND {
		next(s)
		b, f := parseF(s)
		if !b {
			return false, e
		}
		t := (and)([2]exp{e, f})
		return parseEA2(s, t)
	}
	return true, e
}

// N   ::= 0 | 1 | 2
// B   ::= true | false
// V   ::= N | B
func parseV(s *state) (bool, exp) {
	switch {
	case s.tok == ZERO:
		next(s)
		return true, (num)(0)
	case s.tok == ONE:
		next(s)
		return true, (num)(1)
	case s.tok == TWO:
		next(s)
		return true, (num)(2)
	case s.tok == TRUE:
		next(s)
		return true, (boolean)(true)
	case s.tok == FALSE:
		next(s)
		return true, (boolean)(false)
	}

	return false, (err)(0)
}

// F   ::= V | (EP)
func parseF(s *state) (bool, exp) {
	switch {
	case s.tok == OPEN:
		next(s)
		b, e := parseEP(s)
		if !b {
			return false, e
		}
		if s.tok != CLOSE {
			return false, e
		}
		next(s)
		return true, e
	default:
		return parseV(s)
	}
}

func parse(s string) exp {
	st := state{&s, EOS}
	next(&st)
	r, e := parseEP(&st)
	if st.tok == EOS && r {
		return e
	}
	return (err)(0) // dummy value
}

func debug(s string) {
	fmt.Printf("%s", s)
}

func test(s string) {
	e := parse(s)
	fmt.Printf("\nOriginal: %s", s)
	fmt.Printf("\nParsed  : %s", e.pretty())
	fmt.Printf("\n %s", showType(e.infer()))
	fmt.Println("\n========")
}

func testParserGood() {
	fmt.Printf("\n GOOD \n")
	test("1")
	test("1+0")
	test("1 * 2 ")
	test(" (1) ")
	test(" (1 * (2)) ")
	test(" (1 + 2) * 0 ")
	test("true || false")
}

func testParserBad() {
	fmt.Printf("\n BAD \n")
	test("1+")
	test("+ 1")
	test("(((1))")
	test("tru")
	test("fal")
	test("true | false")
	test(" 1 + true")
	test("2 * false")
	test(" 1 || true")
	test(" 1 && false")
}

func main() {
	fmt.Printf("\n")
	testParserGood()
	testParserBad()
}
