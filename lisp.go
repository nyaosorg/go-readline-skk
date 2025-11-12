package skk

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hymkor/sxencode-go/parser"
)

type symbol struct {
	value string
}

type cons struct {
	car any
	cdr any
}

var rawStringToLispString = strings.NewReplacer(
	`\"`, `"`,
	`\\`, `\`,
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

var parser1 = &parser.Parser[any]{
	Cons:   func(car, cdr any) any { return &cons{car: car, cdr: cdr} },
	Number: tryParseAsNumber,
	String: func(s string) any { return rawStringToLispString.Replace(s) },
	Array: func(list []any, dim []int) any {
		array := make([]any, len(list))
		copy(array, list)
		return array
	},
	Keyword: func(s string) any { return s },
	Rune:    func(r rune) any { return r },
	Symbol:  func(s string) any { return symbol{value: s} },
	Null:    func() any { return nil },
	True:    func() any { return true },
}

var rxEscSeq = regexp.MustCompile(`\\[0-9]+`)

func evalSxList(funcs map[string]func([]any) (any, error), sxpr any) (any, error) {
	list := []any{}
	for {
		c, ok := sxpr.(*cons)
		if !ok {
			break
		}
		if cc, ok := c.car.(*cons); ok {
			result, err := evalSxList(funcs, cc)
			if err != nil {
				return nil, err
			}
			list = append(list, result)
		} else {
			list = append(list, c.car)
		}
		sxpr = c.cdr
	}
	if len(list) < 1 {
		return nil, errors.New("too few arguments")
	}
	sym, ok := list[0].(symbol)
	if !ok {
		return nil, errors.New("not a symbol")
	}
	if f, ok := funcs[sym.value]; ok {
		return f(list[1:])
	}
	return nil, errors.New("no such a function")
}

func funConcat(args []any) (any, error) {
	var buffer strings.Builder
	for _, v := range args {
		if s, ok := v.(string); ok {
			buffer.WriteString(s)
		}
	}
	s := buffer.String()
	s = rxEscSeq.ReplaceAllStringFunc(s, func(ss string) string {
		var oct rune = 0
		var b strings.Builder
		for _, c := range ss[1:] {
			oct = (oct * 8) + (c - '0')
		}
		b.WriteRune(oct)
		return b.String()
	})
	return s, nil
}

func funPwd([]any) (any, error) {
	return os.Getwd()
}

func funCurrentTimeString([]any) (any, error) {
	return time.Now().Format(time.ANSIC), nil
}

func funCurrentDate([]any) (any, error) {
	return time.Now().Format("2006年01月02日"), nil
}

func funSubstring(args []any) (any, error) {
	if len(args) != 3 {
		return nil, errors.New("substr: argc error")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, errors.New("substr: not a string")
	}
	start, ok := args[1].(int64)
	if !ok || start < 0 || start >= int64(len(s)) {
		return nil, fmt.Errorf("substr: start index: %v (len=%d)", args[1], len(s))
	}
	end, ok := args[2].(int64)
	if !ok || end < start || end >= int64(len(s)) {
		return nil, fmt.Errorf("substr: end index: %v (len=%d)", args[2], len(s))
	}
	return s[start:end], nil
}

func funSkkVersion(args []any) (any, error) {
	return "go-readline-skk", nil
}

var lispFunctions = map[string]func([]any) (any, error){
	"concat":              funConcat,
	"pwd":                 funPwd,
	"current-time-string": funCurrentTimeString,
	"skk-current-date":    funCurrentDate,
	"substring":           funSubstring,
	"skk-version":         funSkkVersion,
}

func evalSxString(source string) candidateT {
	sxpr, err := parser1.Read(strings.NewReader(source))
	if err != nil {
		return candidateStringT(source)
	}
	return &candidateFuncT{
		source: source,
		f: func() string {
			result, err := evalSxList(lispFunctions, sxpr)
			if err != nil {
				return source
			}
			return fmt.Sprint(result)
		},
	}
}
