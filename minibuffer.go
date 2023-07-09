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
	Recurse(string) MiniBuffer
}

type MiniBufferOnNextLine struct{}

func (MiniBufferOnNextLine) Enter(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\n%s ", prompt)
}

func (MiniBufferOnNextLine) Leave(w io.Writer) (int, error) {
	return io.WriteString(w, "\x1B[F")
}

func (q MiniBufferOnNextLine) Recurse(originalPrompt string) MiniBuffer {
	return &MiniBufferOnCurrentLine{OriginalPrompt: originalPrompt}
}

type MiniBufferOnCurrentLine struct {
	OriginalPrompt string
}

func (q *MiniBufferOnCurrentLine) Enter(w io.Writer, prompt string) (int, error) {
	return fmt.Fprintf(w, "\r%s ", prompt)
}

func (q *MiniBufferOnCurrentLine) Leave(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "\r%s \x1B[K", q.OriginalPrompt)
}

func (q *MiniBufferOnCurrentLine) Recurse(originalPrompt string) MiniBuffer {
	return &MiniBufferOnCurrentLine{OriginalPrompt: originalPrompt}
}

func (M *Mode) message(B *readline.Buffer, text string) {
	M.MiniBuffer.Enter(B.Out, text)
	io.WriteString(B.Out, "\x1B[K")
	M.MiniBuffer.Leave(B.Out)
	B.RepaintAfterPrompt()
}

func (M *Mode) ask1(B *readline.Buffer, prompt string) (string, error) {
	M.MiniBuffer.Enter(B.Out, prompt)
	B.Out.Flush()
	rc, err := B.GetKey()
	io.WriteString(B.Out, "\x1B[2K")
	M.MiniBuffer.Leave(B.Out)
	B.RepaintAfterPrompt()
	return rc, err
}

func (M *Mode) ask(ctx context.Context, B *readline.Buffer, prompt string, ime bool) (string, error) {
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
	if ime {
		m := &Mode{
			User:       M.User,
			System:     M.System,
			MiniBuffer: M.MiniBuffer.Recurse(prompt),
		}
		m.enable(inputNewWord, hiragana)
	}
	defer B.RepaintAfterPrompt()
	return inputNewWord.ReadLine(ctx)
}
