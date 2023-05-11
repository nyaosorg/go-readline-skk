package skk

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/japanese"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

var romajiTable1 = []string{"あ", "い", "う", "え", "お"}

var romajiTable2 = map[string][]string{
	"k": []string{"か", "き", "く", "け", "こ"},
	"s": []string{"さ", "し", "す", "せ", "そ"},
	"t": []string{"た", "ち", "つ", "て", "と"},
	"n": []string{"な", "に", "ぬ", "ね", "の"},
	"h": []string{"は", "ひ", "ふ", "へ", "ほ"},
	"m": []string{"ま", "み", "む", "め", "も"},
	"y": []string{"や", "い", "ゆ", "いぇ", "よ"},
	"r": []string{"ら", "り", "る", "れ", "ろ"},
	"w": []string{"わ", "うぃ", "う", "うぇ", "を"},
	"f": []string{"ふぁ", "ふぃ", "ふ", "ふぇ", "ふぉ"},
	"x": []string{"ぁ", "ぃ", "ぅ", "ぇ", "ぉ"},
	"g": []string{"が", "ぎ", "ぐ", "げ", "ご"},
	"z": []string{"ざ", "じ", "ず", "ぜ", "ぞ"},
	"d": []string{"だ", "ぢ", "づ", "で", "ど"},
	"b": []string{"ば", "び", "ぶ", "べ", "ぼ"},
	"p": []string{"ぱ", "ぴ", "ぷ", "ぺ", "ぽ"},
	"j": []string{"じゃ", "じ", "じゅ", "じぇ", "じょ"},
}

var romajiTable3 = map[string][]string{
	"ky": []string{"きゃ", "きぃ", "きゅ", "きぇ", "きょ"},
	"sh": []string{"しゃ", "し", "しゅ", "しぇ", "しょ"},
	"sy": []string{"しゃ", "しぃ", "しゅ", "しぇ", "しょ"},
	"ty": []string{"ちゃ", "ちぃ", "ちゅ", "ちぇ", "ちょ"},
	"ch": []string{"ちゃ", "ち", "ちゅ", "ちぇ", "ちょ"},
	"ny": []string{"にゃ", "にぃ", "にゅ", "にぇ", "にょ"},
	"hy": []string{"ひゃ", "ひぃ", "ひゅ", "ひぇ", "ひょ"},
	"my": []string{"みゃ", "みぃ", "みゅ", "みぇ", "みょ"},
	"ry": []string{"りゃ", "りぃ", "りゅ", "りぇ", "りょ"},
}

func romajiToKana2char(ctx context.Context, B *rl.Buffer, kana string) rl.Result {
	rl.CmdBackwardDeleteChar.Call(ctx, B)
	return rl.SelfInserter(kana).Call(ctx, B)
}

func romajiToKana3char(ctx context.Context, B *rl.Buffer, kana string) rl.Result {
	rl.CmdBackwardDeleteChar.Call(ctx, B)
	rl.CmdBackwardDeleteChar.Call(ctx, B)
	return rl.SelfInserter(kana).Call(ctx, B)
}

func toString(s io.WriterTo) string {
	var buffer strings.Builder
	s.WriteTo(&buffer)
	return buffer.String()
}

func cmdVowels(ctx context.Context, B *rl.Buffer, aiueo int) rl.Result {
	if B.Cursor >= 2 {
		shiin := B.SubString(B.Cursor-2, B.Cursor)
		if kana, ok := romajiTable3[shiin]; ok {
			return romajiToKana3char(ctx, B, kana[aiueo])
		}
	}
	if B.Cursor >= 1 {
		shiin := toString(B.Buffer[B.Cursor-1].Moji)
		if kana, ok := romajiTable2[shiin]; ok {
			return romajiToKana2char(ctx, B, kana[aiueo])
		}
	}
	return rl.SelfInserter(romajiTable1[aiueo]).Call(ctx, B)
}

func cmdA(ctx context.Context, B *rl.Buffer) rl.Result {
	return cmdVowels(ctx, B, 0)
}

func cmdI(ctx context.Context, B *rl.Buffer) rl.Result {
	return cmdVowels(ctx, B, 1)
}

func cmdU(ctx context.Context, B *rl.Buffer) rl.Result {
	return cmdVowels(ctx, B, 2)
}

func cmdE(ctx context.Context, B *rl.Buffer) rl.Result {
	return cmdVowels(ctx, B, 3)
}

func cmdO(ctx context.Context, B *rl.Buffer) rl.Result {
	return cmdVowels(ctx, B, 4)
}

const (
	markerWhite = "▽"
	markerBlack = "▼"
)

type henkanStart byte

func (h henkanStart) String() string {
	return string(h)
}

func henkanMode(ctx context.Context, B *rl.Buffer, markerPos int, source string, postfix string) rl.Result {
	list, found := jisyo[source]
	if !found {
		// 本来であれば辞書登録モード
		return rl.SelfInserter(" ").Call(ctx, B)
	}
	current := 0
	B.ReplaceAndRepaint(markerPos, markerBlack+list[current]+postfix)
	for {
		B.Out.Flush()
		input, _ := B.GetKey()
		if input < " " {
			removeOne(B, markerPos)
			return rl.CONTINUE
		} else if input == " " {
			current++
			if current >= len(list) {
				current = 0
			}
			B.ReplaceAndRepaint(markerPos, markerBlack+list[current]+postfix)
		} else if input == "x" {
			current--
			if current < 0 {
				B.ReplaceAndRepaint(markerPos, markerWhite+source)
				return rl.CONTINUE
			}
			B.ReplaceAndRepaint(markerPos, markerBlack+list[current]+postfix)
		} else {
			removeOne(B, markerPos)
			return eval(ctx, B, input)
		}
	}
}

func (h henkanStart) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	if markerPos := seekMarker(B); markerPos >= 0 {
		// 送り仮名つき変換
		postfix := string(unicode.ToLower(rune(h)))
		source := B.SubString(markerPos+1, B.Cursor) + postfix
		return henkanMode(ctx, B, markerPos, source, postfix)
	}
	rl.SelfInserter(markerWhite).Call(ctx, B)
	rl.CmdForwardChar.Call(ctx, B)
	switch h {
	case 'a':
		return cmdA(ctx, B)
	case 'i':
		return cmdI(ctx, B)
	case 'u':
		return cmdU(ctx, B)
	case 'e':
		return cmdE(ctx, B)
	case 'o':
		return cmdO(ctx, B)
	}
	return rl.SelfInserter(string(h)).Call(ctx, B)
}

func seekMarker(B *rl.Buffer) int {
	for i := B.Cursor - 1; i >= 0; i-- {
		ch := toString(B.Buffer[i].Moji)
		if ch == markerWhite || ch == markerBlack {
			return i
		}
	}
	return -1
}

var jisyo map[string][]string

func removeOne(B *rl.Buffer, pos int) {
	copy(B.Buffer[pos:], B.Buffer[pos+1:])
	B.Buffer = B.Buffer[:len(B.Buffer)-1]
	B.Cursor--
	B.RepaintAfterPrompt()
}

func cmdHenkan(ctx context.Context, B *rl.Buffer) rl.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return rl.SelfInserter(" ").Call(ctx, B)
	}
	source := B.SubString(markerPos+1, B.Cursor)

	return henkanMode(ctx, B, markerPos, source, "")
}

func eval(ctx context.Context, B *rl.Buffer, input string) rl.Result {
	code := keys.Code(input)
	cmd, ok := B.KeyMap.KeyMap[code]
	if !ok {
		cmd, ok = rl.GlobalKeyMap.KeyMap[code]
		if !ok {
			cmd = rl.SelfInserter(input)
		}
	}
	return cmd.Call(ctx, B)
}

func cmdN(ctx context.Context, B *rl.Buffer) rl.Result {
	rl.SelfInserter("n").Call(ctx, B)
	B.Out.Flush()
	input, _ := B.GetKey()
	switch input {
	case "n":
		rl.CmdBackwardDeleteChar.Call(ctx, B)
		return rl.SelfInserter("ん").Call(ctx, B)
	case "a", "i", "u", "e", "o", "y":
		return eval(ctx, B, input)
	default:
		rl.CmdBackwardDeleteChar.Call(ctx, B)
		rl.SelfInserter("ん").Call(ctx, B)
		return eval(ctx, B, input)
	}
}

type smallTsuChecker string

func (s smallTsuChecker) String() string {
	return "small tsu checker for " + string(s)
}

func (s smallTsuChecker) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	if B.Cursor <= 0 || toString(B.Buffer[B.Cursor-1].Moji) != string(s) {
		return rl.SelfInserter(string(s)).Call(ctx, B)
	}
	rl.CmdBackwardChar.Call(ctx, B)
	rl.SelfInserter("っ").Call(ctx, B)
	return rl.CmdForwardChar.Call(ctx, B)
}

func cmdCtrlJ(ctx context.Context, B *rl.Buffer) rl.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return cmdDisableRomaji(ctx, B)
	}
	// kakutei
	removeOne(B, markerPos)
	return rl.CONTINUE
}

func cmdEnableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	B.BindKey("a", rl.AnonymousCommand(cmdA))
	B.BindKey("i", rl.AnonymousCommand(cmdI))
	B.BindKey("u", rl.AnonymousCommand(cmdU))
	B.BindKey("e", rl.AnonymousCommand(cmdE))
	B.BindKey("o", rl.AnonymousCommand(cmdO))
	B.BindKey("l", rl.AnonymousCommand(cmdDisableRomaji))
	B.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdCtrlJ))
	B.BindKey(" ", rl.AnonymousCommand(cmdHenkan))
	B.BindKey("n", rl.AnonymousCommand(cmdN))

	const upperRomaji = "AIUEOKSTNHMYRWFGZDBPCJ"
	for i, c := range upperRomaji {
		B.BindKey(keys.Code(upperRomaji[i:i+1]), henkanStart(byte(unicode.ToLower(c))))
	}

	const consonantButN = "ksthmyrwfgzdbpcj"
	for i := range consonantButN {
		s := consonantButN[i : i+1]
		B.BindKey(keys.Code(s), smallTsuChecker(s))
	}
	return rl.CONTINUE
}

func cmdDisableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	for i := 'a'; i <= 'z'; i++ {
		s := string(byte(i))
		B.BindKey(keys.Code(s), rl.SelfInserter(s))
	}
	for i := 'A'; i <= 'Z'; i++ {
		s := string(byte(i))
		B.BindKey(keys.Code(s), rl.SelfInserter(s))
	}
	B.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdEnableRomaji))
	return rl.CONTINUE
}

func loadJisyo(filename string) (map[string][]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	decoder := japanese.EUCJP.NewDecoder()

	return readJisyo(decoder.Reader(fd))
}

func readJisyo(r io.Reader) (map[string][]string, error) {
	jisyo := map[string][]string{}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if len(line) <= 0 || line[0] == ';' {
			continue
		}
		source, lists, ok := strings.Cut(line, " /")
		if !ok {
			continue
		}
		values := []string{}
		for {
			one, rest, ok := strings.Cut(lists, "/")
			one, _, _ = strings.Cut(one, ";")
			values = append(values, one)
			if !ok {
				break
			}
			lists = rest
		}
		jisyo[source] = values
	}
	return jisyo, sc.Err()
}

var ErrJisyoNotFound = errors.New("Jisyo not found")

func Setup(jisyoFilenames ...string) error {
	var err error
	for _, fn := range jisyoFilenames {
		jisyo, err = loadJisyo(fn)
		if err == nil {
			rl.GlobalKeyMap.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdEnableRomaji))
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
	}
	return ErrJisyoNotFound
}
