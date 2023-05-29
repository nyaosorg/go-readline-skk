package skk

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"

	rl "github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

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

type _Upper struct {
	H byte
	K *_Kana
	M *Mode
}

func (h *_Upper) String() string {
	return string(h.H)
}

type QueryPrompter interface {
	Prompt(io.Writer, string) (int, error)
	LineFeed(io.Writer) (int, error)
	Recurse(string) QueryPrompter
}

type QueryOnNextLine struct{}

func (_ QueryOnNextLine) Prompt(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\n%s ", prompt)
}

func (_ QueryOnNextLine) LineFeed(w io.Writer) (int, error) {
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

func (M *Mode) ask(ctx context.Context, B *rl.Buffer, prompt string, ime bool) (string, error) {
	B.Out.WriteString("\x1B[?25h")
	B.Out.Flush()
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
	saveMap       map[keys.Code]rl.Command
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
			candidate, _, _ = strings.Cut(list[current], ";")
			B.ReplaceAndRepaint(markerPos, markerBlack+candidate+postfix)
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

func (h *_Upper) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	if markerPos := seekMarker(B); markerPos >= 0 {
		// 送り仮名つき変換
		postfix := string(unicode.ToLower(rune(h.H)))
		source := B.SubString(markerPos+1, B.Cursor) + postfix
		return h.M.henkanMode(ctx, B, markerPos, source, postfix)
	}
	B.InsertAndRepaint(markerWhite)
	switch h.H {
	case 'a':
		return h.K.cmdA(ctx, B)
	case 'i':
		return h.K.cmdI(ctx, B)
	case 'u':
		return h.K.cmdU(ctx, B)
	case 'e':
		return h.K.cmdE(ctx, B)
	case 'o':
		return h.K.cmdO(ctx, B)
	}
	B.InsertAndRepaint(string(h.H))
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
	return "small tsu checker for " + s.post
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

type cmdQ struct {
	mode *Mode
	kana *_Kana
}

func (c *cmdQ) String() string {
	return "SWITCH-KANA"
}

func (c *cmdQ) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	kanaTable[c.kana.switchTo].enableRomaji(B, c.mode)
	return rl.CONTINUE
}

type canBindKey interface {
	BindKey(keys.Code, rl.Command)
	LookupCommand(string) rl.Command
}

func (K *_Kana) enableRomaji(X canBindKey, mode *Mode) {
	X.BindKey("a", rl.AnonymousCommand(K.cmdA))
	X.BindKey("i", rl.AnonymousCommand(K.cmdI))
	X.BindKey("u", rl.AnonymousCommand(K.cmdU))
	X.BindKey("e", rl.AnonymousCommand(K.cmdE))
	X.BindKey("o", rl.AnonymousCommand(K.cmdO))
	X.BindKey("n", rl.AnonymousCommand(K.cmdN))
	X.BindKey(",", rl.SelfInserter("、"))
	X.BindKey(".", rl.SelfInserter("。"))
	X.BindKey("q", &cmdQ{kana: K, mode: mode})
	X.BindKey("-", rl.SelfInserter("ー"))

	const upperRomaji = "AIUEOKSTNHMYRWFGZDBPCJ"
	for i, c := range upperRomaji {
		u := &_Upper{H: byte(unicode.ToLower(c)), K: K, M: mode}
		X.BindKey(keys.Code(upperRomaji[i:i+1]), u)
	}

	const consonantButN = "ksthmyrwfgzdbpcj"
	for i := range consonantButN {
		s := consonantButN[i : i+1]
		X.BindKey(keys.Code(s), &smallTsuChecker{post: s, tsu: K.tsu})
	}
}

func (M *Mode) enableHiragana(X canBindKey) {
	hiragana.enableRomaji(X, M)
	X.BindKey(" ", rl.AnonymousCommand(M.cmdHenkan))
	X.BindKey("l", rl.AnonymousCommand(M.cmdDisableRomaji))
	X.BindKey(keys.CtrlG, rl.AnonymousCommand(M.cmdCtrlG))
	X.BindKey(keys.CtrlJ, rl.AnonymousCommand(M.cmdCtrlJ))
}

func (M *Mode) cmdEnableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	const values = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,.- \x07\n"
	if M.saveMap == nil {
		M.saveMap = map[keys.Code]rl.Command{}
		for i := range values {
			s := values[i : i+1]
			M.saveMap[keys.Code(s)] = B.LookupCommand(s)
		}
	}

	M.enableHiragana(B)
	return rl.CONTINUE
}

func (M *Mode) cmdDisableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	for key, command := range M.saveMap {
		B.BindKey(key, command)
	}
	return rl.CONTINUE
}
