package skk

import (
	"context"
	"fmt"
	"io"

	"github.com/nyaosorg/go-readline-ny"
)

type MiniBuffer interface {
	Enter(io.Writer, string) (int, error)
	Leave(io.Writer) (int, error)
	Recurse() MiniBuffer
}

type MiniBufferOnNextLine struct{}

func (MiniBufferOnNextLine) Enter(w io.Writer, prompt string) (int, error) {
	// この prompt は MiniBufferOnNextLine → MiniBufferOnCurrentLine と呼び出してから戻る際、
	// MiniBufferOnNextLine を再表示する際する時に使う。
	// （ prompt を使わず、別途 io.WriteString(w,prompt) 的な処理を呼び出し元で
	//    やってもらうと、その復元処理ができないかったりする )
	return fmt.Fprintf(w, "\n%s ", prompt)
}

func (MiniBufferOnNextLine) Leave(w io.Writer) (int, error) {
	return io.WriteString(w, "\x1B[F")
}

func (MiniBufferOnNextLine) Recurse() MiniBuffer {
	return MiniBufferOnCurrentLine{}
}

type MiniBufferOnCurrentLine struct {
	OriginalPrompt string
}

func (MiniBufferOnCurrentLine) Enter(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\r%s ", prompt)
}

func (MiniBufferOnCurrentLine) Leave(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "\r\x1B[K")
}

func (MiniBufferOnCurrentLine) Recurse() MiniBuffer {
	return MiniBufferOnCurrentLine{}
}

func (M *Mode) message(B *readline.Buffer, text string) {
	M.MiniBuffer.Enter(B.Out, text)
	io.WriteString(B.Out, "\x1B[K")
	M.MiniBuffer.Leave(B.Out)
	B.RepaintLastLine()
}

func (M *Mode) displayMode(B *readline.Buffer, text string) {
	if _, ok := M.MiniBuffer.(MiniBufferOnNextLine); ok {
		M.message(B, text)
	}
}

func (M *Mode) ask1(B *readline.Buffer, prompt string) (string, error) {
	M.MiniBuffer.Enter(B.Out, prompt)
	B.Out.Flush()
	rc, err := B.GetKey()
	io.WriteString(B.Out, "\x1B[2K")
	M.MiniBuffer.Leave(B.Out)
	B.RepaintLastLine()
	B.Out.Flush()
	return rc, err
}

func (M *Mode) ask(ctx context.Context, B *readline.Buffer, prompt string, ime bool) (string, error) {
	B.Out.Flush()
	inputNewWord := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return M.MiniBuffer.Enter(w, prompt)
		},
		Writer: B.Writer,
		LineFeedWriter: func(_ readline.Result, w io.Writer) (int, error) {
			io.WriteString(w, "\x1B[2K")
			return M.MiniBuffer.Leave(w)
		},
	}
	m := &Mode{
		User:       M.User,
		System:     M.System,
		MiniBuffer: M.MiniBuffer.Recurse(),
		ctrlJ:      M.ctrlJ,
	}
	if ime {
		m.enable(inputNewWord, hiragana)
	} else {
		inputNewWord.BindKey(m.ctrlJ, m)
	}
	m.setupQuitWithLatinMode(inputNewWord)
	inputNewWord.BindKey("\x07", readline.CmdInterrupt)
	rc, err := inputNewWord.ReadLine(ctx)
	B.RepaintLastLine()
	B.Out.Flush()
	return rc, err
}
