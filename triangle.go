package skk

import (
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/go-readline-ny"
)

var (
	isVsCodeTerminal = os.Getenv("VSCODE_PID") != ""

	isWindowsTerminal = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" && !isVsCodeTerminal
)

type triangle rune

func (tr triangle) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c", rune(tr))
	return int64(n), err
}

func (tr triangle) PrintTo(w io.Writer) {
	fmt.Fprintf(w, "%c%c", markerPrefix, rune(tr))
}

func (tr triangle) Width() readline.WidthT {
	if isVsCodeTerminal {
		return 2
	} else if isWindowsTerminal {
		return 1
	} else {
		return 2
	}
}

func insertTriangleAndRepaint(B *readline.Buffer, c triangle) {
	B.Buffer = append(B.Buffer, readline.Cell{})
	copy(B.Buffer[B.Cursor+1:], B.Buffer[B.Cursor:])
	B.Buffer[B.Cursor] = readline.Cell{Moji: c}
	B.Cursor++
	B.DrawFromHead()
}

// ▽は特別扱いなので一旦事前に呼ばれた ReplaceAndRepaint などで
// 普通の文字の「▽」を展開した後
// 後程、改めて、正式版の「▽」に置き換える
func replaceTriangle(B *readline.Buffer, pos int, c triangle) {
	B.Buffer[pos] = readline.Cell{Moji: c}
	B.DrawFromHead()
}
