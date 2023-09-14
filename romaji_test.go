package skk

import (
	"strings"
	"testing"
)

func TestRomaji(t *testing.T) {
	for _, table := range []map[string]string{hiragana.table, katakana.table} {
		for key := range table {
			lastByte := key[0]
			if strings.IndexByte(romajiTrigger, lastByte) < 0 {
				t.Fatalf("%c is not contained in `%s`", lastByte, romajiTrigger)
			}
		}
	}
}
