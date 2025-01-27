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
	whiteMarkerPattern = untilCursor{pattern: regexp.MustCompile(`▽.*?$`)}
	blackMarkerPattern = untilCursor{pattern: regexp.MustCompile(`▼.*?$`)}

	WhiteMarkerHighlight = readline.Highlight{
		Pattern:  whiteMarkerPattern,
		Sequence: "\x1B[0;1;7m", // reverse
	}
	BlackMarkerHighlight = readline.Highlight{
		Pattern:  blackMarkerPattern,
		Sequence: "\x1B[0;1;4m", // underline
	}
)
