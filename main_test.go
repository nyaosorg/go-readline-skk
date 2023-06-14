package skk

import (
	"testing"
)

func TestHanToZen(t *testing.T) {
	list := map[rune]rune{
		'a': 'ａ',
		'z': 'ｚ',
		'A': 'Ａ',
		'Z': 'Ｚ',
		'0': '０',
		'9': '９',
		'!': '！',
		'@': '＠',
		' ': '　',
		'[': '［',
		'|': '｜',
	}

	for source, expect := range list {
		result := hanToZen(source)
		if result != expect {
			t.Fatalf("expect hanToZen('%c')=='%c', but '%c'",
				source, expect, result)
			return
		}
	}
}
