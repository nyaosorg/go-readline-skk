package sxencode

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaosorg/go-readline-skk/internal/sxencode-go/parser"
)

var (
	rxFloat1     = regexp.MustCompile(`^[\-\+]?[0-9]+\.[0-9]+([eE][\-\+]?\d+)?$`)
	rxFloat2     = regexp.MustCompile(`^[\-\+]?[0-9]+[eE][-+]?\d+$`)
	rxInteger    = regexp.MustCompile(`^[\-\+]?[0-9]+$`)
	rxHexInteger = regexp.MustCompile(`^#[Xx][\+\-]?[0-9A-Fa-f]+$`)
	rxOctInteger = regexp.MustCompile(`^#[Oo][\+\-]?[0-7]+$`)
	rxBinInteger = regexp.MustCompile(`^#[Bb][\+\-]?[01]+$`)
)

func tryParseAsFloat(token string) (any, bool, error) {
	if !rxFloat1.MatchString(token) && !rxFloat2.MatchString(token) {
		return nil, false, nil
	}
	val, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, true, err
	}
	return val, true, nil
}

func tryParseAsInt(token string) (any, bool, error) {
	var val int64
	var err error
	if rxInteger.MatchString(token) {
		val, err = strconv.ParseInt(token, 10, 64)
	} else if rxHexInteger.MatchString(token) {
		val, err = strconv.ParseInt(token[2:], 16, 64)
	} else if rxOctInteger.MatchString(token) {
		val, err = strconv.ParseInt(token[2:], 8, 64)
	} else if rxBinInteger.MatchString(token) {
		val, err = strconv.ParseInt(token[2:], 2, 64)
	} else {
		return nil, false, nil
	}
	return val, true, err
}

func tryParseAsNumber(token string) (any, bool, error) {
	if val, ok, err := tryParseAsInt(token); ok {
		return val, true, err
	}
	if val, ok, err := tryParseAsFloat(token); ok {
		return val, true, err
	}
	return nil, false, nil
}

var rawStringToLispString = strings.NewReplacer(
	`\\`, `\`,
	`\"`, `"`,
)

type symbolT struct {
	Value string
}

type consT struct {
	Car any
	Cdr any
}

var parser1 = &parser.Parser[any]{
	Cons:   func(car, cdr any) any { return &consT{Car: car, Cdr: cdr} },
	Number: tryParseAsNumber,
	String: func(s string) any { return rawStringToLispString.Replace(s) },
	Array: func(list []any, dim []int) any {
		array := make([]any, len(list))
		for i, v := range list {
			array[i] = v
		}
		return array
	},
	Keyword: func(s string) any { return s },
	Rune:    func(r rune) any { return r },
	Symbol:  func(s string) any { return symbolT{Value: s} },
	Null:    func() any { return nil },
	True:    func() any { return true },
}
