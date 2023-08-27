package skk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	// "github.com/nyaosorg/go-windows-dbg"
)

func debug(text string) {
	// dbg.Println(text)
}

const (
	markerWhite = "▽"
	markerBlack = "▼"

	msgHiragana = "[か]"
	msgKatakana = "[カ]"
	msgLatin    = ""
	msgAbbrev   = "[aあ]"
	msg0208     = "[英]"
)

type _Trigger struct {
	Key byte
	M   *Mode
}

func (trig *_Trigger) String() string {
	return "SKK_HENKAN_TRIGGER_" + string(trig.Key)
}

// Mode is an instance of SKK. It contains system dictionaries and user dictionaries.
type Mode struct {
	User          *Jisyo
	System        *Jisyo
	MiniBuffer    MiniBuffer
	saveMap       []readline.Command
	kana          *_Kana
	userJisyoPath string
	ctrlJ         keys.Code
}

var rxNumber = regexp.MustCompile(`[0-9]+`)

var rxToNumber = regexp.MustCompile(`#[0123459]`)

var kansuji = map[rune]string{
	'0': "〇",
	'1': "一",
	'2': "二",
	'3': "三",
	'4': "四",
	'5': "五",
	'6': "六",
	'7': "七",
	'8': "八",
	'9': "九",
}

func numberToKanji(s string) string {
	var buffer strings.Builder
	for _, r := range s {
		buffer.WriteString(kansuji[r])
	}
	return buffer.String()
}

func hanToZenString(s string) string {
	var buffer strings.Builder
	for _, r := range s {
		buffer.WriteRune(hanToZen(r))
	}
	return buffer.String()
}

func (M *Mode) _lookup(source string, okuri bool) ([]string, bool) {
	list, ok := M.User.lookup(source, okuri)
	if ok {
		return list, true
	}
	list, ok = M.System.lookup(source, okuri)
	return list, ok
}

func (M *Mode) lookup(source string, okuri bool) ([]string, bool) {
	list, ok := M._lookup(source, okuri)
	if ok {
		return list, ok
	}
	loc := rxNumber.FindStringIndex(source)
	if loc == nil {
		return nil, false
	}
	number := source[loc[0]:loc[1]]
	source = source[:loc[0]] + "#" + source[loc[1]:]
	list, ok = M._lookup(source, okuri)
	if !ok {
		return nil, false
	}
	newList := make([]string, 0, len(list))
	for _, s := range list {
		tmp := rxToNumber.ReplaceAllStringFunc(s, func(ss string) string {
			switch ss[1] {
			case '0': // 無変換
				return number
			case '1': // 全角化
				return hanToZenString(number)
			case '2': // 漢数字で位取りあり
				return numberToKanji(number)
			case '3': // 漢数字で位取りなし
				return numberToKanji(number) // あとでやる
			default:
				return number
			}
		})
		newList = append(newList, tmp)
	}
	return newList, true
}

func unshift(list []string, value string) []string {
	list = append(list, "")
	copy(list[1:], list)
	list[0] = value
	return list
}

func (M *Mode) newCandidate(ctx context.Context, B *readline.Buffer, source string, okuri bool) (string, bool) {
	newWord, err := M.ask(ctx, B, source, true)
	B.RepaintAfterPrompt()
	if err != nil || len(newWord) <= 0 {
		return "", false
	}
	list, _ := M.lookup(source, okuri)

	// 二重登録よけ
	for _, candidate := range list {
		if candidate == newWord {
			return newWord, true
		}
	}
	// リストの先頭に挿入
	M.User.store(source, okuri, unshift(list, newWord))
	return newWord, true
}

const listingStartIndex = 4

func moveTop(list []string, current int) {
	newTop := list[current]
	copy(list[1:current+1], list[:current])
	list[0] = newTop
}

func (M *Mode) henkanMode(ctx context.Context, B *readline.Buffer, markerPos int, source string, postfix string) readline.Result {
	okuri := postfix != ""
	list, found := M.lookup(source, okuri)
	if !found {
		// 辞書登録モード
		result, ok := M.newCandidate(ctx, B, source, okuri)
		if ok {
			// 新変換文字列を展開する
			B.ReplaceAndRepaint(markerPos, result)
			return readline.CONTINUE
		} else {
			// 変換前に一旦戻す
			B.ReplaceAndRepaint(markerPos, markerWhite+source)
			return readline.CONTINUE
		}
	}
	current := 0
	candidate, _, _ := strings.Cut(list[current], ";")
	B.ReplaceAndRepaint(markerPos, markerBlack+candidate+postfix)
	for {
		input, _ := B.GetKey()
		if input == string(keys.CtrlG) {
			B.ReplaceAndRepaint(markerPos, markerWhite+source)
			return readline.CONTINUE
		} else if input < " " {
			if len(postfix) > 0 && postfix[0] == '*' {
				removeOne(B, B.Cursor-len(postfix))
			}
			removeOne(B, markerPos)
			if current > 0 {
				moveTop(list, current)
				M.User.store(source, okuri, list)
			}
			return readline.CONTINUE
		} else if input == " " {
			current++
			if current >= len(list) {
				// 辞書登録モード
				result, ok := M.newCandidate(ctx, B, source, okuri)
				if ok {
					// 新変換文字列を展開する
					B.ReplaceAndRepaint(markerPos, result)
					return readline.CONTINUE
				} else {
					// 変換前に一旦戻す
					B.ReplaceAndRepaint(markerPos, markerWhite+source)
					return readline.CONTINUE
				}
			}
			if current >= listingStartIndex {
				for {
					var buffer strings.Builder
					_current := current
					for _, key := range "ASDFJKL" {
						if _current >= len(list) {
							break
						}
						candidate, _, _ = strings.Cut(list[_current], ";")
						fmt.Fprintf(&buffer, "%c:%s ", key, candidate)
						_current++
					}
					fmt.Fprintf(&buffer, "[残り %d]", len(list)-_current)
					key, err := M.ask1(B, buffer.String())
					if err == nil {
						if index := strings.Index("asdfjkl", key); index >= 0 && current+index < len(list) {
							candidate, _, _ = strings.Cut(list[current+index], ";")
							B.ReplaceAndRepaint(markerPos, candidate)
							return readline.CONTINUE
						} else if key == " " {
							current = _current
							if current >= len(list) {
								// 辞書登録モード
								result, ok := M.newCandidate(ctx, B, source, okuri)
								if ok {
									// 新変換文字列を展開する
									B.ReplaceAndRepaint(markerPos, result)
									return readline.CONTINUE
								} else {
									// 変換前に一旦戻す
									B.ReplaceAndRepaint(markerPos, markerWhite+source)
									return readline.CONTINUE
								}
							}
						} else if key == "x" {
							current -= len("ASDFJKL")
							if current < listingStartIndex {
								if current < 0 {
									current = 0
								}
								break
							}
						} else if key == string(keys.CtrlG) {
							B.ReplaceAndRepaint(markerPos, markerWhite+source)
							return readline.CONTINUE
						}
					}
				}
			} else {
				candidate, _, _ = strings.Cut(list[current], ";")
				B.ReplaceAndRepaint(markerPos, markerBlack+candidate+postfix)
			}
		} else if input == "x" {
			current--
			if current < 0 {
				B.ReplaceAndRepaint(markerPos, markerWhite+source)
				return readline.CONTINUE
			}
			candidate, _, _ = strings.Cut(list[current], ";")
			B.ReplaceAndRepaint(markerPos, markerBlack+candidate+postfix)
		} else if input == "X" {
			prompt := fmt.Sprintf(`really purge "%s /%s/ "?(yes or no)`, source, list[current])
			ans, err := M.ask(ctx, B, prompt, false)
			if err == nil {
				if ans == "y" || ans == "yes" {
					// 本当はシステム辞書を参照しないようLisp構文を
					// セットしなければいけないが、そこまではしない.
					if len(list) <= 1 {
						M.User.remove(source, okuri)
					} else {
						if current+1 < len(list) {
							copy(list[current:], list[current+1:])
						}
						list = list[:len(list)-1]
						M.User.store(source, okuri, list)
					}
					B.ReplaceAndRepaint(markerPos, "")
					return readline.CONTINUE
				}
			}
		} else {
			if len(postfix) > 0 && postfix[0] == '*' {
				removeOne(B, B.Cursor-len(postfix))
			}
			removeOne(B, markerPos)
			if current > 0 {
				moveTop(list, current)
				M.User.store(source, okuri, list)
			}
			return eval(ctx, B, input)
		}
	}
}

func (trig *_Trigger) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	if markerPos := seekMarker(B); markerPos >= 0 {
		// 送り仮名つき変換
		var source strings.Builder
		source.WriteString(B.SubString(markerPos+1, B.Cursor))
		source.WriteByte(trig.Key)

		var postfix string
		if index := strings.IndexByte("aiueo", trig.Key); index >= 0 {
			postfix = trig.M.kana.table[string(trig.Key)]
		} else {
			B.InsertAndRepaint("*" + string(trig.Key))
			var typed strings.Builder
			typed.WriteString(string(trig.Key))
			for {
				key, _ := B.GetKey()
				typed.WriteString(key)
				if value, ok := trig.M.kana.table[typed.String()]; ok {
					return trig.M.henkanMode(ctx, B, markerPos, source.String(), value)
				}
				if len(key) != 1 || !unicode.IsLower(rune(key[0])) {
					return B.LookupCommand(key).Call(ctx, B)
				}
				B.InsertAndRepaint(key)
			}
		}
		return trig.M.henkanMode(ctx, B, markerPos, source.String(), postfix)
	}
	B.InsertAndRepaint(markerWhite)
	r := &_Romaji{kana: trig.M.kana, last: string(trig.Key)}
	return r.Call(ctx, B)
}

func seekMarker(B *readline.Buffer) int {
	for i := B.Cursor - 1; i >= 0; i-- {
		ch := B.Buffer[i].String()
		if ch == markerWhite || ch == markerBlack {
			return i
		}
	}
	return -1
}

func removeOne(B *readline.Buffer, pos int) {
	copy(B.Buffer[pos:], B.Buffer[pos+1:])
	B.Buffer = B.Buffer[:len(B.Buffer)-1]
	B.Cursor--
	B.RepaintAfterPrompt()
}

func (M *Mode) cmdStartHenkan(ctx context.Context, B *readline.Buffer) readline.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		B.InsertAndRepaint(" ")
		return readline.CONTINUE
	}
	source := B.SubString(markerPos+1, B.Cursor)

	return M.henkanMode(ctx, B, markerPos, source, "")
}

func eval(ctx context.Context, B *readline.Buffer, input string) readline.Result {
	return B.LookupCommand(input).Call(ctx, B)
}

func (M *Mode) cmdKakutei(ctx context.Context, B *readline.Buffer) readline.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return M.cmdLatinMode(ctx, B)
	}
	// kakutei
	removeOne(B, markerPos)
	return readline.CONTINUE
}

func (M *Mode) cmdCancel(ctx context.Context, B *readline.Buffer) readline.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return M.cmdLatinMode(ctx, B)
	}
	B.ReplaceAndRepaint(markerPos, "")
	return readline.CONTINUE
}

func (m *Mode) cmdToggleKana(_ context.Context, B *readline.Buffer) readline.Result {
	m.enable(B, kanaTable[m.kana.switchTo])
	if m.kana.switchTo == 1 {
		m.displayMode(B, msgHiragana)
	} else {
		m.displayMode(B, msgKatakana)
	}
	return readline.CONTINUE
}

func (M *Mode) cmdAbbrevMode(ctx context.Context, B *readline.Buffer) readline.Result {
	if seekMarker(B) >= 0 {
		return readline.CONTINUE
	}
	M.restoreKeyMap(B)
	B.InsertAndRepaint(markerWhite)
	B.BindKey(" ", &readline.GoCommand{
		Name: "SKK_ABBREV_START_HENKAN",
		Func: func(ctx context.Context, B *readline.Buffer) readline.Result {
			rc := M.cmdStartHenkan(ctx, B)
			M.enable(B, hiragana)
			M.displayMode(B, msgHiragana)
			return rc
		},
	})
	M.displayMode(B, msgAbbrev)
	return readline.CONTINUE
}

type canLookup interface {
	Lookup(keys.Code) (readline.Command, bool)
}

type CanBindKey interface {
	BindKey(keys.Code, readline.Command)
}

type canKeyMap interface {
	canLookup
	CanBindKey
}

var rxUnicodeCode = regexp.MustCompile(`(?:U+)?([0-9A-Fa-f]{1,8})`)

func (mode *Mode) cmdCodeMode(ctx context.Context, B *readline.Buffer) readline.Result {
	for {
		codeStr, err := mode.ask(ctx, B, "U+", false)
		if err != nil || codeStr == "" {
			return readline.CONTINUE
		}
		m := rxUnicodeCode.FindStringSubmatch(codeStr)
		if m != nil {
			value, err := strconv.ParseUint(m[1], 16, 32)
			if err == nil {
				B.InsertAndRepaint(string(rune(value)))
				return readline.CONTINUE
			}
		}
	}
}

func (mode *Mode) enable(X canKeyMap, K *_Kana) {
	mode.backupKeyMap(X)
	mode.kana = K
	for i := range romajiTrigger {
		c := romajiTrigger[i : i+1]
		X.BindKey(keys.Code(c), &_Romaji{kana: K, last: c})
	}
	const upperRomaji = "AIUEOKSTNHMYRWFGZDBPCJ"
	for i, c := range upperRomaji {
		u := &_Trigger{Key: byte(unicode.ToLower(c)), M: mode}
		X.BindKey(keys.Code(upperRomaji[i:i+1]), u)
	}
	X.BindKey("\\", &readline.GoCommand{Name: "SKK_CODE_MODE", Func: mode.cmdCodeMode})
	X.BindKey("q", &readline.GoCommand{Name: "SKK_TOGGLE_KANA", Func: mode.cmdToggleKana})
	X.BindKey("/", &readline.GoCommand{Name: "SKK_ABBREV_MODE", Func: mode.cmdAbbrevMode})
	X.BindKey(" ", &readline.GoCommand{Name: "SKK_START_HENKAN", Func: mode.cmdStartHenkan})
	X.BindKey("l", &readline.GoCommand{Name: "SKK_LATIN_MODE", Func: mode.cmdLatinMode})
	X.BindKey("L", &readline.GoCommand{Name: "SKK_JISX0208_LATIN_MODE", Func: mode.cmdJis0208LatinMode})
	X.BindKey(keys.CtrlG, &readline.GoCommand{Name: "SKK_CANCEL", Func: mode.cmdCancel})
	X.BindKey(mode.ctrlJ, &readline.GoCommand{Name: "SKK_KAKUTEI", Func: mode.cmdKakutei})
}

func (M *Mode) backupKeyMap(km canLookup) {
	if M.saveMap != nil {
		return
	}
	debug("backupKeyMap")
	M.saveMap = make([]readline.Command, 0, 0x80)
	for i := '\x00'; i <= '\x80'; i++ {
		key := keys.Code(string(i))
		val, _ := km.Lookup(key)
		M.saveMap = append(M.saveMap, val)
	}
}

func (M *Mode) restoreKeyMap(km CanBindKey) {
	debug("restoreKeyMap")
	for i, command := range M.saveMap {
		km.BindKey(keys.Code(string(rune(i))), command)
	}
}

func (M *Mode) cmdLatinMode(ctx context.Context, B *readline.Buffer) readline.Result {
	debug("cmdLatinMode")
	M.restoreKeyMap(B)
	M.displayMode(B, msgLatin)
	return readline.CONTINUE
}

func (M *Mode) cmdAcceptLineWithLatinMode(ctx context.Context, B *readline.Buffer) readline.Result {
	if M.saveMap != nil {
		M.restoreKeyMap(B)
		M.displayMode(B, msgLatin)
	}
	return readline.ENTER
}

func (M *Mode) cmdIntrruptWithLatinMode(ctx context.Context, B *readline.Buffer) readline.Result {
	if M.saveMap != nil {
		M.restoreKeyMap(B)
		M.displayMode(B, msgLatin)
	}
	return readline.INTR
}

func hanToZen(c rune) rune {
	if c < ' ' || c >= '\x7f' {
		return c
	}
	if c == ' ' {
		return '　'
	}
	return c - ' ' + '\uFF00'
}

func (M *Mode) cmdJis0208LatinMode(ctx context.Context, B *readline.Buffer) readline.Result {
	for i := rune(' '); i < '\x7F'; i++ {
		z := string(hanToZen(i))
		B.BindKey(keys.Code(string(i)), &readline.GoCommand{
			Name: "SKK_JISX0208_LATIN_INSERT_" + z,
			Func: func(_ context.Context, B *readline.Buffer) readline.Result {
				B.InsertAndRepaint(z)
				return readline.CONTINUE
			}})
	}
	B.BindKey(M.ctrlJ, &readline.GoCommand{
		Name: "SKK_JISX0208_LATIN_KAKUTEI",
		Func: func(ctx context.Context, B *readline.Buffer) readline.Result {
			M.restoreKeyMap(B)
			M.enable(B, hiragana)
			M.displayMode(B, msgHiragana)
			return readline.CONTINUE
		},
	})
	M.displayMode(B, msg0208)
	return readline.CONTINUE
}
