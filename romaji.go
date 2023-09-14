package skk

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny"
)

type _Kana struct {
	table    map[string]string
	switchTo int
}

var kanaTable = []*_Kana{
	hiragana,
	katakana,
}

const romajiTrigger = "aiueokstnhmyrwfgzdbpcj',.-[]Qx"

var hiragana = &_Kana{
	table: map[string]string{
		"a": "あ", "i": "い", "u": "う", "e": "え", "o": "お", "'": "'",
		",": "、", ".": "。", "-": "ー", "[": "「", "]": "」", "Q": markerWhite,

		"ka": "か", "ki": "き", "ku": "く", "ke": "け", "ko": "こ", "kk": "っk", "nk": "んk",
		"sa": "さ", "si": "し", "su": "す", "se": "せ", "so": "そ", "ss": "っs", "ns": "んs",
		"ta": "た", "ti": "ち", "tu": "つ", "te": "て", "to": "と", "tt": "っt", "nt": "んt",
		"na": "な", "ni": "に", "nu": "ぬ", "ne": "ね", "no": "の", "nn": "ん", "n'": "ん",
		"ha": "は", "hi": "ひ", "hu": "ふ", "he": "へ", "ho": "ほ", "hh": "っh", "nh": "んh",
		"ma": "ま", "mi": "み", "mu": "む", "me": "め", "mo": "も", "mm": "っm", "nm": "んm",
		"ya": "や", "yu": "ゆ", "yo": "よ", "yy": "っy",
		"ra": "ら", "ri": "り", "ru": "る", "re": "れ", "ro": "ろ", "rr": "っr", "nr": "んr",
		"wa": "わ", "wo": "を", "ww": "っw", "nw": "んw",
		"fa": "ふぁ", "fi": "ふぃ", "fu": "ふ", "fe": "ふぇ", "fo": "ふぉ", "ff": "っf", "nf": "んf",
		"xa": "ぁ", "xi": "ぃ", "xu": "ぅ", "xe": "ぇ", "xo": "ぉ", "xx": "っx", "nx": "んx",
		"ga": "が", "gi": "ぎ", "gu": "ぐ", "ge": "げ", "go": "ご", "ng": "んg",
		"za": "ざ", "zi": "じ", "zu": "ず", "ze": "ぜ", "zo": "ぞ", "zz": "っz", "nz": "んz",
		"da": "だ", "di": "ぢ", "du": "づ", "de": "で", "do": "ど", "dd": "っd", "nd": "んd",
		"ba": "ば", "bi": "び", "bu": "ぶ", "be": "べ", "bo": "ぼ", "bb": "っb", "nb": "んb",
		"pa": "ぱ", "pi": "ぴ", "pu": "ぷ", "pe": "ぺ", "po": "ぽ", "pp": "っp", "np": "んp",
		"ja": "じゃ", "ji": "じ", "ju": "じゅ", "je": "じぇ", "jo": "じょ", "jj": "っj", "nj": "んj",

		"kya": "きゃ", "kyi": "きぃ", "kyu": "きゅ", "kye": "きぇ", "kyo": "きょ",
		"sha": "しゃ", "shi": "し", "shu": "しゅ", "she": "しぇ", "sho": "しょ",
		"sya": "しゃ", "syi": "しぃ", "syu": "しゅ", "sye": "しぇ", "syo": "しょ",
		"tya": "ちゃ", "tyi": "ちぃ", "tyu": "ちゅ", "tye": "ちぇ", "tyo": "ちょ",
		"cha": "ちゃ", "chi": "ち", "chu": "ちゅ", "che": "ちぇ", "cho": "ちょ",
		"nya": "にゃ", "nyi": "にぃ", "nyu": "にゅ", "nye": "にぇ", "nyo": "にょ",
		"hya": "ひゃ", "hyi": "ひぃ", "hyu": "ひゅ", "hye": "ひぇ", "hyo": "ひょ",
		"mya": "みゃ", "myi": "みぃ", "myu": "みゅ", "mye": "みぇ", "myo": "みょ",
		"rya": "りゃ", "ryi": "りぃ", "ryu": "りゅ", "rye": "りぇ", "ryo": "りょ",
		"dha": "でゃ", "dhi": "でぃ", "dhu": "でゅ", "dhe": "でぇ", "dho": "でょ",
		"dya": "ぢゃ", "dyi": "ぢぃ", "dyu": "ぢゅ", "dye": "ぢぇ", "dyo": "ぢょ",
		"gya": "ぎゃ", "gyi": "ぎぃ", "gyu": "ぎゅ", "gye": "ぎぇ", "gyo": "ぎょ",
		"xya": "ゃ", "xyu": "ゅ", "xyo": "ょ", "xtu": "っ",

		"zh": "←", "zj": "↓", "zk": "↑", "zl": "→",
		"xtsu": "っ",
	},
	switchTo: 1,
}

var katakana = &_Kana{
	table: map[string]string{
		"a": "ア", "i": "イ", "u": "ウ", "e": "エ", "o": "オ", "'": "'",
		",": "、", ".": "。", "-": "ー", "[": "「", "]": "」", "Q": markerWhite,

		"ka": "カ", "ki": "キ", "ku": "ク", "ke": "ケ", "ko": "コ", "kk": "ッk", "nk": "ンk",
		"sa": "サ", "si": "シ", "su": "ス", "se": "セ", "so": "ソ", "ss": "ッs", "ns": "ンs",
		"ta": "タ", "ti": "チ", "tu": "ツ", "te": "テ", "to": "ト", "tt": "ッt", "nt": "ンt",
		"na": "ナ", "ni": "ニ", "nu": "ヌ", "ne": "ネ", "no": "ノ", "nn": "ン", "n'": "ン",
		"ha": "ハ", "hi": "ヒ", "hu": "フ", "he": "ヘ", "ho": "ホ", "hh": "ッh", "nh": "ンh",
		"ma": "マ", "mi": "ミ", "mu": "ム", "me": "メ", "mo": "モ", "mm": "ッm", "nm": "ンm",
		"ya": "ヤ", "yu": "ユ", "yo": "ヨ", "yy": "ッy",
		"ra": "ラ", "ri": "リ", "ru": "ル", "re": "レ", "ro": "ロ", "rr": "ッr", "nr": "ンr",
		"wa": "ワ", "wi": "ウィ", "wu": "ウ", "we": "ウェ", "wo": "ヲ", "ww": "ッw", "nw": "ンw",
		"fa": "ファ", "fi": "フィ", "fu": "フ", "fe": "フェ", "fo": "フォ", "ff": "ッf", "nf": "ンf",
		"xa": "ァ", "xi": "ィ", "xu": "ゥ", "xe": "ェ", "xo": "ォ", "xx": "ッx", "nx": "ンx",
		"ga": "ガ", "gi": "ギ", "gu": "グ", "ge": "ゲ", "go": "ゴ", "gg": "ッg", "ng": "ンg",
		"za": "ザ", "zi": "ジ", "zu": "ズ", "ze": "ゼ", "zo": "ゾ", "zz": "ッz", "nz": "ンz",
		"da": "ダ", "di": "ジ", "du": "ヅ", "de": "デ", "do": "ド", "dd": "ッd", "nd": "ンd",
		"ba": "バ", "bi": "ビ", "bu": "ブ", "be": "ベ", "bo": "ボ", "bb": "ッb", "nb": "ンb",
		"pa": "パ", "pi": "ピ", "pu": "プ", "pe": "ペ", "po": "ポ", "pp": "ッp", "np": "ンp",
		"ja": "ジャ", "ji": "ジ", "ju": "ジュ", "je": "ジェ", "jo": "ジョ", "jj": "ッj", "nj": "ンj",

		"kya": "キャ", "kyi": "キ", "kyu": "キュ", "kye": "キェ", "kyo": "キョ",
		"sha": "シャ", "shi": "シ", "shu": "シュ", "she": "シェ", "sho": "ショ",
		"sya": "シャ", "syi": "シィ", "syu": "シュ", "sye": "シェ", "syo": "ショ",
		"tya": "チャ", "tyi": "チィ", "tyu": "チュ", "tye": "チェ", "tyo": "チョ",
		"cha": "チャ", "chi": "チ", "chu": "チュ", "che": "チェ", "cho": "チョ",
		"nya": "ニャ", "nyi": "ニィ", "nyu": "ニュ", "nye": "ニェ", "nyo": "ニョ",
		"hya": "ヒャ", "hyi": "ヒィ", "hyu": "ヒュ", "hye": "ヒェ", "hyo": "ヒョ",
		"mya": "ミャ", "myi": "ミィ", "myu": "ミュ", "mye": "ミェ", "myo": "ミョ",
		"rya": "リャ", "ryi": "リィ", "ryu": "リュ", "rye": "リェ", "ryo": "リョ",
		"dha": "デャ", "dhi": "ディ", "dhu": "デュ", "dhe": "デェ", "dho": "デョ",
		"dya": "ヂャ", "dyi": "ヂィ", "dyu": "ヂュ", "dye": "ヂェ", "dyo": "ヂョ",
		"gya": "ギャ", "gyi": "ギィ", "gyu": "ギュ", "gye": "ギェ", "gyo": "ギョ",
		"xya": "ャ", "xyu": "ュ", "xyo": "ョ", "xtu": "ッ",

		"zh": "←", "zj": "↓", "zk": "↑", "zl": "→",
		"xtsu": "ッ",
	},
	switchTo: 0,
}

type _Romaji struct {
	kana *_Kana
	last string
}

func (R *_Romaji) String() string {
	return "SKK_ROMAJI_" + R.last
}

func (R *_Romaji) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	if value, ok := R.kana.table[R.last]; ok {
		B.InsertAndRepaint(value)
		return readline.CONTINUE
	}
	var buffer strings.Builder
	buffer.WriteString(R.last)
	from := B.Cursor
	B.InsertAndRepaint(string(R.last))
	for {
		input, _ := B.GetKey()
		if len(input) != 1 {
			eval(ctx, B, input)
			return readline.CONTINUE
		}
		if !unicode.IsLetter(rune(input[0])) {
			B.ReplaceAndRepaint(from, buffer.String())
			return readline.CONTINUE
		}
		c := unicode.ToLower(rune(input[0]))
		buffer.WriteRune(c)
		if value, ok := R.kana.table[buffer.String()]; ok {
			B.ReplaceAndRepaint(from, value)
			if u, _ := utf8.DecodeLastRuneInString(value); u != c {
				return readline.CONTINUE
			}
			buffer.Reset()
			buffer.WriteRune(c)
			from = B.Cursor - 1
		} else {
			B.InsertAndRepaint(string(c))
		}
	}
}
