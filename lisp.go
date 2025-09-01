package skk

import (
	"math/big"
	"os"
	"regexp"
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

var parser1 = &parser.Parser[any]{
	Cons:   func(car, cdr any) any { return &cons{car: car, cdr: cdr} },
	Int:    func(n int64) any { return n },
	BigInt: func(n *big.Int) any { return n },
	Float:  func(f float64) any { return f },
	String: func(s string) any { return s },
	Array: func(list []any, dim []int) any {
		array := make([]any, len(list))
		for i, v := range list {
			array[i] = v
		}
		return array
	},
	Keyword: func(s string) any { return s },
	Rune:    func(r rune) any { return r },
	Symbol:  func(s string) any { return symbol{value: s} },
	Null:    func() any { return nil },
	True:    func() any { return true },
}

var rxEscSeq = regexp.MustCompile(`\\[0-9]+`)

func evalSxString(source string) candidateT {
	sxpr, err := parser1.Read(strings.NewReader(source))
	if err != nil {
		return candidateStringT(source)
	}
	list := []any{}
	for {
		c, ok := sxpr.(*cons)
		if !ok {
			break
		}
		list = append(list, c.car)
		sxpr = c.cdr
	}
	if len(list) < 1 {
		return candidateStringT(source)
	}
	sym, ok := list[0].(symbol)
	if !ok {
		return candidateStringT(source)
	}
	switch sym.value {
	case "concat":
		var buffer strings.Builder
		for _, v := range list[1:] {
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
		return &candidateFuncT{
			source: source,
			f:      func() string { return s },
		}
	case "pwd": // `/pwd`
		return &candidateFuncT{
			source: source,
			f: func() string {
				wd, err := os.Getwd()
				if err != nil {
					return source
				}
				return wd
			},
		}
	case "current-time-string": // `/now`, `/time`
		return &candidateFuncT{
			source: source,
			f: func() string {
				return time.Now().Format(time.ANSIC)
			},
		}
	case "skk-current-date":
		return &candidateFuncT{
			source: source,
			f: func() string {
				return time.Now().Format("2006年01月02日")
			},
		}
	}
	return candidateStringT(source)
}
