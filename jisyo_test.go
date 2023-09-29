package skk

import (
	"io"
	"strings"
	"testing"
)

func TestPeekLine(t *testing.T) {
	sample := "hogehoge\n" +
		"ahaha\n" +
		"ihihi\n" +
		"ohoho\n" +
		"fehehe\n" +
		"ufufu"

	var r io.Reader = strings.NewReader(sample)
	r, line, err := peekLine(r)

	if err != nil {
		t.Fatalf("ERR=%s", err.Error())
	}
	if line != "hogehoge\n" {
		t.Fatalf("expect %s, but %s", "hogehoge\n", line)
	}
	all, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll: %s", err.Error())
	}
	if string(all) != sample {
		t.Fatalf("io.ReadAll: expect `%s` but `%s`", sample, string(all))
	}
}
