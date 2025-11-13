package skk

import (
	"regexp"

	"github.com/nyaosorg/go-readline-ny"
)

type untilCursor struct {
	pattern *regexp.Regexp
}

func (uc untilCursor) FindAllStringIndex(str string, n int) [][]int {
	return uc.pattern.FindAllStringIndex(str[:-n-1], n)
}

var (
	triangleOutlinePattern = untilCursor{pattern: regexp.MustCompile(`▽.*?$`)}
	triangleFilledPattern  = untilCursor{pattern: regexp.MustCompile(`▼.*?$`)}

	TriangleOutlineHighlight = readline.Highlight{
		Pattern:  triangleOutlinePattern,
		Sequence: "\x1B[0;1;7m", // reverse
	}
	WhiteMarkerHighlight = TriangleOutlineHighlight

	TriangleFilledHighlight = readline.Highlight{
		Pattern:  triangleFilledPattern,
		Sequence: "\x1B[0;1;4m", // underline
	}
	BlackMarkerHighlight = TriangleFilledHighlight
)
