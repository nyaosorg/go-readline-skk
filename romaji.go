package skk

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny"
)

type _Kana struct {
	table        map[string]string
	hiraKataSwTo int
	hanzenSwTo   int
	modeStr      string
}

func (K *_Kana) Query(romaji string) (string, bool) {
	romaji = strings.ToLower(romaji)
	if result, ok := K.table[romaji]; ok {
		return result, true
	}
	if len(romaji) >= 2 && romaji[0] == romaji[1] {
		if result, ok := K.table[romaji[1:]]; ok {
			return K.table["xtsu"] + result, true
		}
	}
	return "", false
}

var kanaTable = []*_Kana{
	hiragana,
	katakana,
	hankakuHiragana,
	hankakuKatakana,
}

const romajiTrigger = "aiueokstnhmyrwfgzdbpcj',.-[]Qx"

var hiragana = &_Kana{
	table: map[string]string{
		"a": "あ", "i": "い", "u": "う", "e": "え", "o": "お", "'": "'",
		",": "、", ".": "。", "-": "ー", "[": "「", "]": "」",

		"ka": "か", "ki": "き", "ku": "く", "ke": "け", "ko": "こ", "nk": "んk",
		"sa": "さ", "si": "し", "su": "す", "se": "せ", "so": "そ", "ns": "んs",
		"ta": "た", "ti": "ち", "tu": "つ", "te": "て", "to": "と", "nt": "んt",
		"na": "な", "ni": "に", "nu": "ぬ", "ne": "ね", "no": "の", "nn": "ん", "n'": "ん",
		"ha": "は", "hi": "ひ", "hu": "ふ", "he": "へ", "ho": "ほ", "nh": "んh",
		"ma": "ま", "mi": "み", "mu": "む", "me": "め", "mo": "も", "nm": "んm",
		"ya": "や", "yu": "ゆ", "yo": "よ", "yy": "っy",
		"ra": "ら", "ri": "り", "ru": "る", "re": "れ", "ro": "ろ", "nr": "んr",
		"wa": "わ", "wo": "を", "ww": "っw", "nw": "んw",
		"fa": "ふぁ", "fi": "ふぃ", "fu": "ふ", "fe": "ふぇ", "fo": "ふぉ", "nf": "んf",
		"xa": "ぁ", "xi": "ぃ", "xu": "ぅ", "xe": "ぇ", "xo": "ぉ", "nx": "んx",
		"ga": "が", "gi": "ぎ", "gu": "ぐ", "ge": "げ", "go": "ご", "ng": "んg",
		"za": "ざ", "zi": "じ", "zu": "ず", "ze": "ぜ", "zo": "ぞ", "nz": "んz",
		"da": "だ", "di": "ぢ", "du": "づ", "de": "で", "do": "ど", "nd": "んd",
		"ba": "ば", "bi": "び", "bu": "ぶ", "be": "べ", "bo": "ぼ", "nb": "んb",
		"pa": "ぱ", "pi": "ぴ", "pu": "ぷ", "pe": "ぺ", "po": "ぽ", "np": "んp",
		"ja": "じゃ", "ji": "じ", "ju": "じゅ", "je": "じぇ", "jo": "じょ", "nj": "んj",

		"kya": "きゃ", "kyi": "きぃ", "kyu": "きゅ", "kye": "きぇ", "kyo": "きょ",
		"sha": "しゃ", "shi": "し", "shu": "しゅ", "she": "しぇ", "sho": "しょ",
		"sya": "しゃ", "syi": "しぃ", "syu": "しゅ", "sye": "しぇ", "syo": "しょ",
		"tha": "てぁ", "thi": "てぃ", "thu": "てゅ", "the": "てぇ", "tho": "てょ",
		"tya": "ちゃ", "tyi": "ちぃ", "tyu": "ちゅ", "tye": "ちぇ", "tyo": "ちょ",
		"cha": "ちゃ", "chi": "ち", "chu": "ちゅ", "che": "ちぇ", "cho": "ちょ",
		"nya": "にゃ", "nyi": "にぃ", "nyu": "にゅ", "nye": "にぇ", "nyo": "にょ",
		"hya": "ひゃ", "hyi": "ひぃ", "hyu": "ひゅ", "hye": "ひぇ", "hyo": "ひょ",
		"mya": "みゃ", "myi": "みぃ", "myu": "みゅ", "mye": "みぇ", "myo": "みょ",
		"rya": "りゃ", "ryi": "りぃ", "ryu": "りゅ", "rye": "りぇ", "ryo": "りょ",
		"dha": "でゃ", "dhi": "でぃ", "dhu": "でゅ", "dhe": "でぇ", "dho": "でょ",
		"dya": "ぢゃ", "dyi": "ぢぃ", "dyu": "ぢゅ", "dye": "ぢぇ", "dyo": "ぢょ",
		"gya": "ぎゃ", "gyi": "ぎぃ", "gyu": "ぎゅ", "gye": "ぎぇ", "gyo": "ぎょ",
		"bya": "びゃ", "byi": "びぃ", "byu": "びゅ", "bye": "びぇ", "byo": "びょ",
		"pya": "ぴゃ", "pyi": "ぴぃ", "pyu": "ぴゅ", "pye": "ぴぇ", "pyo": "ぴょ",
		"xya": "ゃ", "xyu": "ゅ", "xyo": "ょ", "xtu": "っ",

		"zh": "←", "zj": "↓", "zk": "↑", "zl": "→",
		"z,": "‥", "z-": "～", "z.": "…", "z/": "・", "z[": "『", "z]": "』",

		"z1": "○", "z2": "▽", "z3": "△", "z4": "□", "z5": "◇",
		"z6": "☆", "z7": "◎", "z8": "〔", "z9": "〕", "z0": "∞",
		"z^": "※", "z\\": "￥", "z@": "〃", "z;": "゛", "z:": "゜",

		"z!": "●", "z\"": "▼", "z#": "▲", "z$": "■ ", "z%": "◆",
		"z&": "★", "z'": "♪", "z(": "【", "z)": "】", "z=": "≒",
		"z~": "≠", "z|": "〒", "z`": "“", "z+": "±", "z*": "×",
		"z<": "≦", "z>": "≧", "z?": "÷", "z_": "―", "z ": "\u3000",

		"xtsu": "っ",
	},
	hiraKataSwTo: 1,
	hanzenSwTo:   2,
	modeStr:      msgHiragana,
}

var katakana = &_Kana{
	table: map[string]string{
		"a": "ア", "i": "イ", "u": "ウ", "e": "エ", "o": "オ", "'": "'",
		",": "、", ".": "。", "-": "ー", "[": "「", "]": "」",

		"ka": "カ", "ki": "キ", "ku": "ク", "ke": "ケ", "ko": "コ", "nk": "ンk",
		"sa": "サ", "si": "シ", "su": "ス", "se": "セ", "so": "ソ", "ns": "ンs",
		"ta": "タ", "ti": "チ", "tu": "ツ", "te": "テ", "to": "ト", "nt": "ンt",
		"na": "ナ", "ni": "ニ", "nu": "ヌ", "ne": "ネ", "no": "ノ", "nn": "ン", "n'": "ン",
		"ha": "ハ", "hi": "ヒ", "hu": "フ", "he": "ヘ", "ho": "ホ", "nh": "ンh",
		"ma": "マ", "mi": "ミ", "mu": "ム", "me": "メ", "mo": "モ", "nm": "ンm",
		"ya": "ヤ", "yu": "ユ", "yo": "ヨ",
		"ra": "ラ", "ri": "リ", "ru": "ル", "re": "レ", "ro": "ロ", "nr": "ンr",
		"wa": "ワ", "wi": "ウィ", "wu": "ウ", "we": "ウェ", "wo": "ヲ", "nw": "ンw",
		"fa": "ファ", "fi": "フィ", "fu": "フ", "fe": "フェ", "fo": "フォ", "nf": "ンf",
		"xa": "ァ", "xi": "ィ", "xu": "ゥ", "xe": "ェ", "xo": "ォ", "nx": "ンx",
		"ga": "ガ", "gi": "ギ", "gu": "グ", "ge": "ゲ", "go": "ゴ", "ng": "ンg",
		"za": "ザ", "zi": "ジ", "zu": "ズ", "ze": "ゼ", "zo": "ゾ", "nz": "ンz",
		"da": "ダ", "di": "ジ", "du": "ヅ", "de": "デ", "do": "ド", "nd": "ンd",
		"ba": "バ", "bi": "ビ", "bu": "ブ", "be": "ベ", "bo": "ボ", "nb": "ンb",
		"pa": "パ", "pi": "ピ", "pu": "プ", "pe": "ペ", "po": "ポ", "np": "ンp",
		"ja": "ジャ", "ji": "ジ", "ju": "ジュ", "je": "ジェ", "jo": "ジョ", "nj": "ンj",

		"kya": "キャ", "kyi": "キ", "kyu": "キュ", "kye": "キェ", "kyo": "キョ",
		"sha": "シャ", "shi": "シ", "shu": "シュ", "she": "シェ", "sho": "ショ",
		"sya": "シャ", "syi": "シィ", "syu": "シュ", "sye": "シェ", "syo": "ショ",
		"tha": "テァ", "thi": "ティ", "thu": "テュ", "the": "テェ", "tho": "テョ",
		"tya": "チャ", "tyi": "チィ", "tyu": "チュ", "tye": "チェ", "tyo": "チョ",
		"cha": "チャ", "chi": "チ", "chu": "チュ", "che": "チェ", "cho": "チョ",
		"nya": "ニャ", "nyi": "ニィ", "nyu": "ニュ", "nye": "ニェ", "nyo": "ニョ",
		"hya": "ヒャ", "hyi": "ヒィ", "hyu": "ヒュ", "hye": "ヒェ", "hyo": "ヒョ",
		"mya": "ミャ", "myi": "ミィ", "myu": "ミュ", "mye": "ミェ", "myo": "ミョ",
		"rya": "リャ", "ryi": "リィ", "ryu": "リュ", "rye": "リェ", "ryo": "リョ",
		"dha": "デャ", "dhi": "ディ", "dhu": "デュ", "dhe": "デェ", "dho": "デョ",
		"dya": "ヂャ", "dyi": "ヂィ", "dyu": "ヂュ", "dye": "ヂェ", "dyo": "ヂョ",
		"gya": "ギャ", "gyi": "ギィ", "gyu": "ギュ", "gye": "ギェ", "gyo": "ギョ",
		"bya": "ビャ", "byi": "ビィ", "byu": "ビュ", "bye": "ビェ", "byo": "ビョ",
		"pya": "ピャ", "pyi": "ピィ", "pyu": "ピュ", "pye": "ピェ", "pyo": "ピョ",
		"xya": "ャ", "xyu": "ュ", "xyo": "ョ", "xtu": "ッ",

		"zh": "←", "zj": "↓", "zk": "↑", "zl": "→",
		"z,": "‥", "z-": "～", "z.": "…", "z/": "・", "z[": "『", "z]": "』",

		"z1": "○", "z2": "▽", "z3": "△", "z4": "□", "z5": "◇",
		"z6": "☆", "z7": "◎", "z8": "〔", "z9": "〕", "z0": "∞",
		"z^": "※", "z\\": "￥", "z@": "〃", "z;": "゛", "z:": "゜",

		"z!": "●", "z\"": "▼", "z#": "▲", "z$": "■ ", "z%": "◆",
		"z&": "★", "z'": "♪", "z(": "【", "z)": "】", "z=": "≒",
		"z~": "≠", "z|": "〒", "z`": "“", "z+": "±", "z*": "×",
		"z<": "≦", "z>": "≧", "z?": "÷", "z_": "―", "z ": "\u3000",

		"xtsu": "ッ",
	},
	hiraKataSwTo: 0,
	hanzenSwTo:   3,
	modeStr:      msgKatakana,
}

var hankaku = map[string]string{
	"a": "ｱ", "i": "ｲ", "u": "ｳ", "e": "ｴ", "o": "ｵ", "'": "'",
	",": "､", ".": "｡", "-": "ｰ", "[": "｢", "]": "｣",

	"ka": "ｶ", "ki": "ｷ", "ku": "ｸ", "ke": "ｹ", "ko": "ｺ", "nk": "ﾝk",
	"sa": "ｻ", "si": "ｼ", "su": "ｽ", "se": "ｾ", "so": "ｿ", "ns": "ﾝs",
	"ta": "ﾀ", "ti": "ﾁ", "tu": "ﾂ", "te": "ﾃ", "to": "ﾄ", "nt": "ﾝt",
	"na": "ﾅ", "ni": "ﾆ", "nu": "ﾇ", "ne": "ﾈ", "no": "ﾉ", "nn": "ﾝ", "n'": "ﾝ",
	"ha": "ﾊ", "hi": "ﾋ", "hu": "ﾌ", "he": "ﾍ", "ho": "ﾎ", "nh": "ﾝh",
	"ma": "ﾏ", "mi": "ﾐ", "mu": "ﾑ", "me": "ﾒ", "mo": "ﾓ", "nm": "ﾝm",
	"ya": "ﾔ", "yu": "ﾕ", "yo": "ﾖ",
	"ra": "ﾗ", "ri": "ﾘ", "ru": "ﾙ", "re": "ﾚ", "ro": "ﾛ", "nr": "ﾝr",
	"wa": "ﾜ", "wi": "ｳｨ", "wu": "ｳ", "we": "ｳｪ", "wo": "ｦ", "nw": "ﾝw",
	"fa": "ﾌｧ", "fi": "ﾌｨ", "fu": "ﾌ", "fe": "ﾌｪ", "fo": "ﾌｫ", "nf": "ﾝf",
	"xa": "ｧ", "xi": "ｨ", "xu": "ｩ", "xe": "ｪ", "xo": "ｫ", "nx": "ﾝx",
	"ga": "ｶﾞ", "gi": "ｷﾞ", "gu": "ｸﾞ", "ge": "ｹﾞ", "go": "ｺﾞ", "ng": "ﾝg",
	"za": "ｻﾞ", "zi": "ｼﾞ", "zu": "ｽﾞ", "ze": "ｾﾞ", "zo": "ｿﾞ", "nz": "ﾝz",
	"da": "ﾀﾞ", "di": "ｼﾞ", "du": "ﾂﾞ", "de": "ﾃﾞ", "do": "ﾄﾞ", "nd": "ﾝd",
	"ba": "ﾊﾞ", "bi": "ﾋﾞ", "bu": "ﾌﾞ", "be": "ﾍﾞ", "bo": "ﾎﾞ", "nb": "ﾝb",
	"pa": "ﾊﾟ", "pi": "ﾋﾟ", "pu": "ﾌﾟ", "pe": "ﾍﾟ", "po": "ﾎﾟ", "np": "ﾝp",
	"ja": "ｼﾞｬ", "ji": "ｼﾞ", "ju": "ｼﾞｭ", "je": "ｼﾞｪ", "jo": "ｼﾞｮ", "nj": "ﾝj",

	"kya": "ｷｬ", "kyi": "ｷ", "kyu": "ｷｭ", "kye": "ｷｪ", "kyo": "ｷｮ",
	"sha": "ｼｬ", "shi": "ｼ", "shu": "ｼｭ", "she": "ｼｪ", "sho": "ｼｮ",
	"sya": "ｼｬ", "syi": "ｼｨ", "syu": "ｼｭ", "sye": "ｼｪ", "syo": "ｼｮ",
	"tha": "ﾃｧ", "thi": "ﾃｨ", "thu": "ﾃｭ", "the": "ﾃｪ", "tho": "ﾃｮ",
	"tya": "ﾁｬ", "tyi": "ﾁｨ", "tyu": "ﾁｭ", "tye": "ﾁｪ", "tyo": "ﾁｮ",
	"cha": "ﾁｬ", "chi": "ﾁ", "chu": "ﾁｭ", "che": "ﾁｪ", "cho": "ﾁｮ",
	"nya": "ﾆｬ", "nyi": "ﾆｨ", "nyu": "ﾆｭ", "nye": "ﾆｪ", "nyo": "ﾆｮ",
	"hya": "ﾋｬ", "hyi": "ﾋｨ", "hyu": "ﾋｭ", "hye": "ﾋｪ", "hyo": "ﾋｮ",
	"mya": "ﾐｬ", "myi": "ﾐｨ", "myu": "ﾐｭ", "mye": "ﾐｪ", "myo": "ﾐｮ",
	"rya": "ﾘｬ", "ryi": "ﾘｨ", "ryu": "ﾘｭ", "rye": "ﾘｪ", "ryo": "ﾘｮ",
	"dha": "ﾃﾞｬ", "dhi": "ﾃﾞｨ", "dhu": "ﾃﾞｭ", "dhe": "ﾃﾞｪ", "dho": "ﾃﾞｮ",
	"dya": "ﾁﾞｬ", "dyi": "ﾁﾞｨ", "dyu": "ﾁﾞｭ", "dye": "ﾁﾞｪ", "dyo": "ﾁﾞｮ",
	"gya": "ｷﾞｬ", "gyi": "ｷﾞｨ", "gyu": "ｷﾞｭ", "gye": "ｷﾞｪ", "gyo": "ｷﾞｮ",
	"bya": "ﾋﾞｬ", "byi": "ﾋﾞｨ", "byu": "ﾋﾞｭ", "bye": "ﾋﾞｪ", "byo": "ﾋﾞｮ",
	"pya": "ﾋﾟｬ", "pyi": "ﾋﾟｨ", "pyu": "ﾋﾟｭ", "pye": "ﾋﾟｪ", "pyo": "ﾋﾟｮ",
	"xya": "ｬ", "xyu": "ｭ", "xyo": "ｮ", "xtu": "ｯ",

	"zh": "←", "zj": "↓", "zk": "↑", "zl": "→",
	"z,": "‥", "z-": "～", "z.": "…", "z/": "・", "z[": "『", "z]": "』",

	"z1": "○", "z2": "▽", "z3": "△", "z4": "□", "z5": "◇",
	"z6": "☆", "z7": "◎", "z8": "〔", "z9": "〕", "z0": "∞",
	"z^": "※", "z\\": "￥", "z@": "〃", "z;": "゛", "z:": "゜",

	"z!": "●", "z\"": "▼", "z#": "▲", "z$": "■ ", "z%": "◆",
	"z&": "★", "z'": "♪", "z(": "【", "z)": "】", "z=": "≒",
	"z~": "≠", "z|": "〒", "z`": "“", "z+": "±", "z*": "×",
	"z<": "≦", "z>": "≧", "z?": "÷", "z_": "―", "z ": " ",

	"xtsu": "ｯ",
}

var hankakuHiragana = &_Kana{
	table:        hankaku,
	hiraKataSwTo: 1,
	hanzenSwTo:   0,
	modeStr:      msgHankaku,
}

var hankakuKatakana = &_Kana{
	table:        hankaku,
	hiraKataSwTo: 1,
	hanzenSwTo:   1,
	modeStr:      msgHankaku,
}

type _Romaji struct {
	kana *_Kana
	last string
}

func (R *_Romaji) String() string {
	return "SKK_ROMAJI_" + R.last
}

func (R *_Romaji) Call(ctx context.Context, B *readline.Buffer) readline.Result {
	if value, ok := R.kana.Query(R.last); ok {
		B.InsertAndRepaint(value)
		return readline.CONTINUE
	}
	var buffer strings.Builder
	buffer.WriteString(R.last)
	from := B.Cursor
	B.InsertAndRepaint(string(R.last))
	for {
		input, _ := B.GetKey()
		if len(input) != 1 || input[0] < ' ' {
			eval(ctx, B, input)
			return readline.CONTINUE
		}
		c := unicode.ToLower(rune(input[0]))
		buffer.WriteRune(c)
		if value, ok := R.kana.Query(buffer.String()); ok {
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
