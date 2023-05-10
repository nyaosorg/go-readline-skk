package skk

import (
	"context"
	"strings"
	"unicode"

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
	"w": []string{"わ", "うぃ", "う", "うぇ", "を"},
	"f": []string{"ふぁ", "ふぃ", "ふ", "ふぇ", "ふぉ"},
	"x": []string{"ぁ", "ぃ", "ぅ", "ぇ", "ぉ"},
}

var romajiTable3 = map[string][]string{
	"ky": []string{"きゃ", "きぃ", "きゅ", "きぇ", "きょ"},
	"sy": []string{"しゃ", "しぃ", "しゅ", "しぇ", "しょ"},
	"ty": []string{"ちゃ", "ちぃ", "ちゅ", "ちぇ", "ちょ"},
	"ny": []string{"にゃ", "にぃ", "にゅ", "にぇ", "にょ"},
	"hy": []string{"ひゃ", "ひぃ", "ひゅ", "ひぇ", "ひょ"},
	"my": []string{"みゃ", "みぃ", "みゅ", "みぇ", "みょ"},
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

func cmdVowels(ctx context.Context, B *rl.Buffer, aiueo int) rl.Result {
	if B.Cursor >= 2 {
		var buffer strings.Builder
		B.Buffer[B.Cursor-2].Moji.WriteTo(&buffer)
		B.Buffer[B.Cursor-1].Moji.WriteTo(&buffer)
		shiin := buffer.String()
		if kana, ok := romajiTable3[shiin]; ok {
			return romajiToKana3char(ctx, B, kana[aiueo])
		}
	}
	if B.Cursor >= 1 {
		var buffer strings.Builder
		B.Buffer[B.Cursor-1].Moji.WriteTo(&buffer)
		shiin := buffer.String()
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

const henkanMarker = "▽"

type henkanStart byte

func (h henkanStart) String() string {
	return string(h)
}

func (h henkanStart) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	rl.SelfInserter(henkanMarker).Call(ctx, B)
	rl.CmdForwardChar.Call(ctx, B)
	switch h {
	case 'a':
		return cmdA(ctx, B)
	case 'i':
		return cmdA(ctx, B)
	case 'u':
		return cmdA(ctx, B)
	case 'e':
		return cmdA(ctx, B)
	case 'o':
		return cmdA(ctx, B)
	}
	return rl.SelfInserter(string(h)).Call(ctx, B)
}

func cmdEnableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	B.BindKey("a", rl.AnonymousCommand(cmdA))
	B.BindKey("i", rl.AnonymousCommand(cmdI))
	B.BindKey("u", rl.AnonymousCommand(cmdU))
	B.BindKey("e", rl.AnonymousCommand(cmdE))
	B.BindKey("o", rl.AnonymousCommand(cmdO))
	B.BindKey("l", rl.AnonymousCommand(cmdDisableRomaji))
	B.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdDisableRomaji))

	for _, c := range "AIUEOKSTNHMYRWF" {
		B.BindKey(keys.Code(string(c)), henkanStart(byte(unicode.ToLower(c))))
	}
	return rl.CONTINUE
}

func cmdDisableRomaji(ctx context.Context, B *rl.Buffer) rl.Result {
	for _, s := range "aiueolAIUEOKSTNHMYRWF" {
		B.BindKey(keys.Code(string(s)), rl.SelfInserter(string(s)))
	}
	B.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdEnableRomaji))
	return rl.CONTINUE
}

func init() {
	rl.GlobalKeyMap.BindKey(keys.CtrlJ, rl.AnonymousCommand(cmdEnableRomaji))
}
