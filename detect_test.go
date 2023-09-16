package skk

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestDetectAndRewind(t *testing.T) {
	expects := []string{
		"hogehoge\n",
		"ahaha\n",
		"ihihi\n",
		"ohoho\n",
		"fehehe\n",
		"ufufu",
	}
	sample := strings.Join(expects, "")

	for _, target := range expects {
		r, err := detectAndRewind(strings.NewReader(sample), func(line string) bool {
			return line == target
		})
		if err != nil {
			t.Fatal(err.Error())
		}
		br := bufio.NewReader(r)

		for i, expected := range expects {
			line, err := br.ReadString('\n')
			if err != nil {
				if i != len(expects)-1 || err != io.EOF {
					t.Fatal(err.Error())
				}
			}
			if line != expected {
				t.Fatalf("expected %s but %s", expected, line)
			}
		}
	}
}
