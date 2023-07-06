package skk

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-windows-dbg"
)

func debug(text string) {
	dbg.Println(text)
}

type _Kana struct {
	table1   []string
	table2   map[string][]string
	table3   map[string][]string
	table4   map[string][]string
	n        string
	tsu      string
	switchTo int
}

var kanaTable = []*_Kana{
	hiragana,
	katakana,
}

var hiragana = &_Kana{
	table1: []string{"あ", "い", "う", "え", "お"},
	table2: map[string][]string{
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
	},
	table3: map[string][]string{
		"ky": []string{"きゃ", "きぃ", "きゅ", "きぇ", "きょ"},
		"sh": []string{"しゃ", "し", "しゅ", "しぇ", "しょ"},
		"sy": []string{"しゃ", "しぃ", "しゅ", "しぇ", "しょ"},
		"ty": []string{"ちゃ", "ちぃ", "ちゅ", "ちぇ", "ちょ"},
		"ch": []string{"ちゃ", "ち", "ちゅ", "ちぇ", "ちょ"},
		"ny": []string{"にゃ", "にぃ", "にゅ", "にぇ", "にょ"},
		"hy": []string{"ひゃ", "ひぃ", "ひゅ", "ひぇ", "ひょ"},
		"my": []string{"みゃ", "みぃ", "みゅ", "みぇ", "みょ"},
		"ry": []string{"りゃ", "りぃ", "りゅ", "りぇ", "りょ"},
		"xt": []string{"っぁ", "っぃ", "っ", "っぇ", "っぉ"},
		"dh": []string{"でゃ", "でぃ", "でゅ", "でぇ", "でょ"},
		"dy": []string{"ぢゃ", "ぢぃ", "ぢゅ", "ぢぇ", "ぢょ"},
		"gy": []string{"ぎゃ", "ぎぃ", "ぎゅ", "ぎぇ", "ぎょ"},
	},
	table4: map[string][]string{
		"xts": []string{"っぁ", "っぃ", "っ", "っぇ", "っぉ"},
	},
	n:        "ん",
	tsu:      "っ",
	switchTo: 1,
}

var katakana = &_Kana{
	table1: []string{"ア", "イ", "ウ", "エ", "オ"},
	table2: map[string][]string{
		"k": []string{"カ", "キ", "ク", "ケ", "コ"},
		"s": []string{"サ", "シ", "ス", "セ", "ソ"},
		"t": []string{"タ", "チ", "ツ", "テ", "ト"},
		"n": []string{"ナ", "ニ", "ヌ", "ネ", "ノ"},
		"h": []string{"ハ", "ヒ", "フ", "ヘ", "ホ"},
		"m": []string{"マ", "ミ", "ム", "メ", "モ"},
		"y": []string{"ヤ", "イ", "ユ", "イェ", "ヨ"},
		"r": []string{"ラ", "リ", "ル", "レ", "ロ"},
		"w": []string{"ワ", "ウィ", "ウ", "ウェ", "ヲ"},
		"f": []string{"ファ", "フィ", "フ", "フェ", "フォ"},
		"x": []string{"ァ", "ィ", "ゥ", "ェ", "ォ"},
		"g": []string{"ガ", "ギ", "グ", "ゲ", "ゴ"},
		"z": []string{"ザ", "ジ", "ズ", "ゼ", "ゾ"},
		"d": []string{"ダ", "ジ", "ヅ", "デ", "ド"},
		"b": []string{"バ", "ビ", "ブ", "ベ", "ボ"},
		"p": []string{"パ", "ピ", "プ", "ペ", "ポ"},
		"j": []string{"ジャ", "ジ", "ジュ", "ジェ", "ジョ"},
	},
	table3: map[string][]string{
		"ky": []string{"キャ", "キ", "キュ", "キェ", "キョ"},
		"sh": []string{"シャ", "シ", "シュ", "シェ", "ショ"},
		"sy": []string{"シャ", "シィ", "シュ", "シェ", "ショ"},
		"ty": []string{"チャ", "チィ", "チュ", "チェ", "チョ"},
		"ch": []string{"チャ", "チ", "チュ", "チェ", "チョ"},
		"ny": []string{"ニャ", "ニィ", "ニュ", "ニェ", "ニョ"},
		"hy": []string{"ヒャ", "ヒィ", "ヒュ", "ヒェ", "ヒョ"},
		"my": []string{"ミャ", "ミ", "ミュ", "ミェ", "ミョ"},
		"ry": []string{"リャ", "リィ", "リュ", "リェ", "リョ"},
		"xt": []string{"ッァ", "ッィ", "ッ", "ッェ", "ッォ"},
		"dh": []string{"デャ", "ディ", "デュ", "デェ", "デョ"},
		"dy": []string{"ヂャ", "ヂィ", "ヂュ", "ヂェ", "ヂョ"},
		"gy": []string{"ギャ", "ギィ", "ギュ", "ギェ", "ギョ"},
	},
	table4: map[string][]string{
		"xts": []string{"ッァ", "ッィ", "ッ", "ッェ", "ッォ"},
	},
	n:        "ン",
	tsu:      "ッ",
	switchTo: 0,
}

func (K *_Kana) cmdVowels(ctx context.Context, B *rl.Buffer, aiueo int) rl.Result {
	if B.Cursor >= 3 {
		shiin := B.SubString(B.Cursor-3, B.Cursor)
		if kana, ok := K.table4[shiin]; ok {
			B.ReplaceAndRepaint(B.Cursor-3, kana[aiueo])
			return rl.CONTINUE
		}
	}
	if B.Cursor >= 2 {
		shiin := B.SubString(B.Cursor-2, B.Cursor)
		if kana, ok := K.table3[shiin]; ok {
			B.ReplaceAndRepaint(B.Cursor-2, kana[aiueo])
			return rl.CONTINUE
		}
	}
	if B.Cursor >= 1 {
		shiin := B.Buffer[B.Cursor-1].String()
		if kana, ok := K.table2[shiin]; ok {
			B.ReplaceAndRepaint(B.Cursor-1, kana[aiueo])
			return rl.CONTINUE
		}
	}
	B.InsertAndRepaint(K.table1[aiueo])
	return rl.CONTINUE
}

func (K *_Kana) cmdA(ctx context.Context, B *rl.Buffer) rl.Result {
	return K.cmdVowels(ctx, B, 0)
}

func (K *_Kana) cmdI(ctx context.Context, B *rl.Buffer) rl.Result {
	return K.cmdVowels(ctx, B, 1)
}

func (K *_Kana) cmdU(ctx context.Context, B *rl.Buffer) rl.Result {
	return K.cmdVowels(ctx, B, 2)
}

func (K *_Kana) cmdE(ctx context.Context, B *rl.Buffer) rl.Result {
	return K.cmdVowels(ctx, B, 3)
}

func (K *_Kana) cmdO(ctx context.Context, B *rl.Buffer) rl.Result {
	return K.cmdVowels(ctx, B, 4)
}

const (
	markerWhite = "▽"
	markerBlack = "▼"
)

type _Trigger struct {
	Key byte
	M   *Mode
}

func (trig *_Trigger) String() string {
	return "SKK_HENKAN_TRIGGER_" + string(trig.Key)
}

type QueryPrompter interface {
	Prompt(io.Writer, string) (int, error)
	LineFeed(io.Writer) (int, error)
	Recurse(string) QueryPrompter
}

type QueryOnNextLine struct{}

func (QueryOnNextLine) Prompt(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\n%s ", prompt)
}

func (QueryOnNextLine) LineFeed(w io.Writer) (int, error) {
	return io.WriteString(w, "\r\x1B[K\x1B[A")
}

func (q QueryOnNextLine) Recurse(originalPrompt string) QueryPrompter {
	return &QueryOnCurrentLine{OriginalPrompt: originalPrompt}
}

type QueryOnCurrentLine struct {
	OriginalPrompt string
}

func (q *QueryOnCurrentLine) Prompt(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\r%s ", prompt)
}

func (q *QueryOnCurrentLine) LineFeed(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "\r%s \x1B[K", q.OriginalPrompt)
}

func (q *QueryOnCurrentLine) Recurse(originalPrompt string) QueryPrompter {
	return &QueryOnCurrentLine{OriginalPrompt: originalPrompt}
}

func (M *Mode) ask1(B *rl.Buffer, prompt string) (string, error) {
	M.QueryPrompter.Prompt(B.Out, prompt)
	B.Out.Flush()
	rc, err := B.GetKey()
	M.QueryPrompter.LineFeed(B.Out)
	B.RepaintAfterPrompt()
	return rc, err
}

func (M *Mode) ask(ctx context.Context, B *rl.Buffer, prompt string, ime bool) (string, error) {
	inputNewWord := &rl.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return M.QueryPrompter.Prompt(w, prompt)
		},
		Writer: B.Writer,
		LineFeedWriter: func(_ rl.Result, w io.Writer) (int, error) {
			return M.QueryPrompter.LineFeed(w)
		},
	}
	if ime {
		m := &Mode{
			User:          M.User,
			System:        M.System,
			QueryPrompter: M.QueryPrompter.Recurse(prompt),
		}
		m.enableHiragana(inputNewWord)
	}
	defer B.RepaintAfterPrompt()
	return inputNewWord.ReadLine(ctx)
}

// Mode is an instance of SKK. It contains system dictionaries and user dictionaries.
type Mode struct {
	User          Jisyo
	System        Jisyo
	QueryPrompter QueryPrompter
	saveMap       []rl.Command
	kana          *_Kana
}

func (M *Mode) newCandidate(ctx context.Context, B *rl.Buffer, source string) (string, bool) {
	newWord, err := M.ask(ctx, B, source, true)
	B.RepaintAfterPrompt()
	if err != nil || len(newWord) <= 0 {
		return "", false
	}
	list, ok := M.User[source]
	if !ok {
		list = M.System[source]
	}
	// 二重登録よけ
	for _, candidate := range list {
		if candidate == newWord {
			return newWord, true
		}
	}
	// リストの先頭に挿入
	list = append(list, "")
	copy(list[1:], list)
	list[0] = newWord
	M.User[source] = list
	return newWord, true
}

const listingStartIndex = 4

func (M *Mode) henkanMode(ctx context.Context, B *rl.Buffer, markerPos int, source string, postfix string) rl.Result {
	list, found := M.User[source]
	if !found {
		list, found = M.System[source]
	}
	if !found {
		// 辞書登録モード
		result, ok := M.newCandidate(ctx, B, source)
		if ok {
			// 新変換文字列を展開する
			B.ReplaceAndRepaint(markerPos, result)
			return rl.CONTINUE
		} else {
			// 変換前に一旦戻す
			B.ReplaceAndRepaint(markerPos, markerWhite+source)
			return rl.CONTINUE
		}
	}
	current := 0
	candidate, _, _ := strings.Cut(list[current], ";")
	B.ReplaceAndRepaint(markerPos, markerBlack+candidate+postfix)
	for {
		input, _ := B.GetKey()
		if input == string(keys.CtrlG) {
			B.ReplaceAndRepaint(markerPos, markerWhite+source)
			return rl.CONTINUE
		} else if input < " " {
			removeOne(B, markerPos)
			return rl.CONTINUE
		} else if input == " " {
			current++
			if current >= len(list) {
				// 辞書登録モード
				result, ok := M.newCandidate(ctx, B, source)
				if ok {
					// 新変換文字列を展開する
					B.ReplaceAndRepaint(markerPos, result)
					return rl.CONTINUE
				} else {
					// 変換前に一旦戻す
					B.ReplaceAndRepaint(markerPos, markerWhite+source)
					return rl.CONTINUE
				}
			}
			if current >= listingStartIndex {
				for {
					var buffer strings.Builder
					_current := current
					for _, key := range "ASDFJKL:" {
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
						if index := strings.Index("asdfjkl:", key); index >= 0 {
							candidate, _, _ = strings.Cut(list[current+index], ";")
							B.ReplaceAndRepaint(markerPos, candidate)
							return rl.CONTINUE
						} else if key == " " {
							current = _current
						} else if key == "x" {
							current -= len("ASDFJKL:")
							if current < listingStartIndex {
								if current < 0 {
									current = 0
								}
								break
							}
						} else if key == string(keys.CtrlG) {
							B.ReplaceAndRepaint(markerPos, markerWhite+source)
							return rl.CONTINUE
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
				return rl.CONTINUE
			}
			B.ReplaceAndRepaint(markerPos, markerBlack+list[current]+postfix)
		} else if input == "X" {
			prompt := fmt.Sprintf(`really purse "%s /%s/ "?(yes or no)`, source, list[current])
			ans, err := M.ask(ctx, B, prompt, false)
			if err == nil {
				if ans == "y" || ans == "yes" {
					// 本当はシステム辞書を参照しないようLisp構文を
					// セットしなければいけないが、そこまではしない.
					if len(list) <= 1 {
						delete(M.User, source)
					} else {
						if current+1 < len(list) {
							copy(list[current:], list[current+1:])
						}
						list = list[:len(list)-1]
						M.User[source] = list
					}
					B.ReplaceAndRepaint(markerPos, "")
					return rl.CONTINUE
				}
			}
		} else {
			removeOne(B, markerPos)
			return eval(ctx, B, input)
		}
	}
}

func (trig *_Trigger) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	if markerPos := seekMarker(B); markerPos >= 0 {
		// 送り仮名つき変換
		var source strings.Builder
		source.WriteString(B.SubString(markerPos+1, B.Cursor))
		source.WriteByte(trig.Key)

		var postfix string
		if index := strings.IndexByte("aiueo", trig.Key); index >= 0 {
			postfix = trig.M.kana.table1[index]
		} else {
			postfix = string(trig.Key)
		}
		return trig.M.henkanMode(ctx, B, markerPos, source.String(), postfix)
	}
	B.InsertAndRepaint(markerWhite)
	switch trig.Key {
	case 'a':
		return trig.M.kana.cmdA(ctx, B)
	case 'i':
		return trig.M.kana.cmdI(ctx, B)
	case 'u':
		return trig.M.kana.cmdU(ctx, B)
	case 'e':
		return trig.M.kana.cmdE(ctx, B)
	case 'o':
		return trig.M.kana.cmdO(ctx, B)
	}
	B.InsertAndRepaint(string(trig.Key))
	return rl.CONTINUE
}

func seekMarker(B *rl.Buffer) int {
	for i := B.Cursor - 1; i >= 0; i-- {
		ch := B.Buffer[i].String()
		if ch == markerWhite || ch == markerBlack {
			return i
		}
	}
	return -1
}

func removeOne(B *rl.Buffer, pos int) {
	copy(B.Buffer[pos:], B.Buffer[pos+1:])
	B.Buffer = B.Buffer[:len(B.Buffer)-1]
	B.Cursor--
	B.RepaintAfterPrompt()
}

func (M *Mode) cmdHenkan(ctx context.Context, B *rl.Buffer) rl.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		B.InsertAndRepaint(" ")
		return rl.CONTINUE
	}
	source := B.SubString(markerPos+1, B.Cursor)

	return M.henkanMode(ctx, B, markerPos, source, "")
}

func eval(ctx context.Context, B *rl.Buffer, input string) rl.Result {
	return B.LookupCommand(input).Call(ctx, B)
}

func (K *_Kana) cmdN(ctx context.Context, B *rl.Buffer) rl.Result {
	B.InsertAndRepaint("n")
	input, _ := B.GetKey()
	switch input {
	case "n":
		B.ReplaceAndRepaint(B.Cursor-1, K.n)
		return rl.CONTINUE
	case "a", "i", "u", "e", "o", "y":
		return eval(ctx, B, input)
	default:
		B.ReplaceAndRepaint(B.Cursor-1, K.n)
		return eval(ctx, B, input)
	}
}

type smallTsuChecker struct {
	post string
	tsu  string
}

func (s *smallTsuChecker) String() string {
	return "SKK_SMALL_TSU_CHCKER_FOR_" + s.post
}

func (s *smallTsuChecker) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	if B.Cursor <= 0 || B.Buffer[B.Cursor-1].String() != s.post {
		B.InsertAndRepaint(s.post)
	} else {
		B.ReplaceAndRepaint(B.Cursor-1, s.tsu+B.Buffer[B.Cursor-1].String())
	}
	return rl.CONTINUE
}

func (M *Mode) cmdCtrlJ(ctx context.Context, B *rl.Buffer) rl.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return M.cmdDisableRomaji(ctx, B)
	}
	// kakutei
	removeOne(B, markerPos)
	return rl.CONTINUE
}

func (M *Mode) cmdCtrlG(ctx context.Context, B *rl.Buffer) rl.Result {
	markerPos := seekMarker(B)
	if markerPos < 0 {
		return M.cmdDisableRomaji(ctx, B)
	}
	B.ReplaceAndRepaint(markerPos, "")
	return rl.CONTINUE
}

func (m *Mode) cmdQ(_ context.Context, B *rl.Buffer) rl.Result {
	m.kana = kanaTable[m.kana.switchTo]
	m.kana.enableRomaji(B, m)
	return rl.CONTINUE
}

func (M *Mode) cmdSlash(ctx context.Context, B *rl.Buffer) rl.Result {
	if seekMarker(B) >= 0 {
		return rl.CONTINUE
	}
	M.restoreKeyMap(&B.KeyMap)
	B.InsertAndRepaint(markerWhite)
	B.BindKey(" ", &rl.GoCommand{
		Name: "SKK_SPACE_AFTER_SLASH",
		Func: func(ctx context.Context, B *rl.Buffer) rl.Result {
			rc := M.cmdHenkan(ctx, B)
			M.kana.enableRomaji(B, M)
			return rc
		},
	})
	return rl.CONTINUE
}

type canBindKey interface {
	BindKey(keys.Code, rl.Command)
	LookupCommand(string) rl.Command
}

func (K *_Kana) enableRomaji(X canBindKey, mode *Mode) {
	X.BindKey("a", &rl.GoCommand{Name: "SKK_A", Func: K.cmdA})
	X.BindKey("i", &rl.GoCommand{Name: "SKK_I", Func: K.cmdI})
	X.BindKey("u", &rl.GoCommand{Name: "SKK_U", Func: K.cmdU})
	X.BindKey("e", &rl.GoCommand{Name: "SKK_E", Func: K.cmdE})
	X.BindKey("o", &rl.GoCommand{Name: "SKK_O", Func: K.cmdO})
	X.BindKey("n", &rl.GoCommand{Name: "SKK_N", Func: K.cmdN})
	X.BindKey(",", rl.SelfInserter("、"))
	X.BindKey(".", rl.SelfInserter("。"))
	X.BindKey("q", &rl.GoCommand{Name: "SKK_Q", Func: mode.cmdQ})
	X.BindKey("-", rl.SelfInserter("ー"))
	X.BindKey("[", rl.SelfInserter("「"))
	X.BindKey("]", rl.SelfInserter("」"))
	X.BindKey("/", &rl.GoCommand{Name: "SKK_SLASH", Func: mode.cmdSlash})

	const upperRomaji = "AIUEOKSTNHMYRWFGZDBPCJ"
	for i, c := range upperRomaji {
		u := &_Trigger{Key: byte(unicode.ToLower(c)), M: mode}
		X.BindKey(keys.Code(upperRomaji[i:i+1]), u)
	}

	const consonantButN = "ksthmyrwfgzdbpcj"
	for i := range consonantButN {
		s := consonantButN[i : i+1]
		X.BindKey(keys.Code(s), &smallTsuChecker{post: s, tsu: K.tsu})
	}
}

func (M *Mode) enableHiragana(X canBindKey) {
	debug("enableHiragana")
	M.kana = hiragana
	hiragana.enableRomaji(X, M)
	X.BindKey(" ", &rl.GoCommand{Name: "SKK_SPACE", Func: M.cmdHenkan})
	X.BindKey("l", &rl.GoCommand{Name: "SKK_L", Func: M.cmdDisableRomaji})
	X.BindKey("L", &rl.GoCommand{Name: "SKK_LARGE_L", Func: M.largeL})
	X.BindKey(keys.CtrlG, &rl.GoCommand{Name: "SKK_CTRL_G", Func: M.cmdCtrlG})
	X.BindKey(keys.CtrlJ, &rl.GoCommand{Name: "SKK_CTRL_J", Func: M.cmdCtrlJ})
}

func (M *Mode) backupKeyMap(km *rl.KeyMap) {
	if M.saveMap != nil {
		return
	}
	debug("backupKeyMap")
	M.saveMap = make([]rl.Command, 0, 0x80)
	for i := '\x00'; i <= '\x80'; i++ {
		key := keys.Code(string(i))
		val, _ := km.Lookup(key)
		M.saveMap = append(M.saveMap, val)
	}
}

func (M *Mode) restoreKeyMap(km *rl.KeyMap) {
	debug("restoreKeyMap")
	for i, command := range M.saveMap {
		km.BindKey(keys.Code(string(rune(i))), command)
	}
}

func (M *Mode) cmdEnableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	debug("cmdEnableRomaji")
	M.backupKeyMap(&B.KeyMap)
	M.enableHiragana(B)
	return rl.CONTINUE
}

func (M *Mode) cmdDisableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	debug("cmdDisableRomaji")
	M.restoreKeyMap(&B.KeyMap)
	return rl.CONTINUE
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

func (M *Mode) largeL(ctx context.Context, B *rl.Buffer) rl.Result {
	for i := rune(' '); i < '\x7F'; i++ {
		z := string(hanToZen(i))
		B.BindKey(keys.Code(string(i)), &rl.GoCommand{
			Name: "SKK_INSERT_" + z,
			Func: func(_ context.Context, B *rl.Buffer) rl.Result {
				B.InsertAndRepaint(z)
				return rl.CONTINUE
			}})
	}
	B.BindKey(keys.CtrlJ, &rl.GoCommand{
		Name: "SKK_CTRL_J_ON_LARGE_L",
		Func: func(ctx context.Context, B *rl.Buffer) rl.Result {
			M.restoreKeyMap(&B.Editor.KeyMap)
			return M.cmdEnableRomaji(ctx, B)
		},
	})
	return rl.CONTINUE
}
