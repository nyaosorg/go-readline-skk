package skk

import (
	"bufio"
	"bytes"
	"io"
)

func detectAndRewind(r io.Reader, f func(string) bool) (io.Reader, error) {
	var buffer bytes.Buffer
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		buffer.WriteString(line)
		if f(line) || err != nil {
			io.CopyN(&buffer, br, int64(br.Buffered()))
			if err == io.EOF {
				return &buffer, nil
			}
			return io.MultiReader(&buffer, r), err
		}
	}
}
