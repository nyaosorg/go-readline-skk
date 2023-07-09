package skk

import (
	"fmt"
	"io"
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
	return io.WriteString(w, "\r\x1B[K\x1B[A")
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
