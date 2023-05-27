package skk

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
)

// Jisyo is a dictionary that contains user or system dictionary.
type Jisyo map[string][]string

// Load reads the contents of an dictionary from a file as EUC-JP.
func (j Jisyo) Load(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return j.ReadEucJp(fd)
}

// Load reads the contents of an dictionary from io.Reader as EUC-JP
func (j Jisyo) ReadEucJp(r io.Reader) error {
	decoder := japanese.EUCJP.NewDecoder()
	return j.Read(decoder.Reader(r))
}

// Load reads the contents of an dictionary from io.Reader as UTF8
func (j Jisyo) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if len(line) <= 0 || line[0] == ';' {
			continue
		}
		source, lists, ok := strings.Cut(line, " /")
		if !ok {
			continue
		}
		values := []string{}
		for {
			one, rest, ok := strings.Cut(lists, "/")
			if one != "" {
				values = append(values, one)
			}
			if !ok {
				break
			}
			lists = rest
		}
		j[source] = values
	}
	return sc.Err()
}

type writeCounter struct {
	n   int64
	err error
}

func (w *writeCounter) Try(n int, err error) bool {
	w.n += int64(n)
	w.err = err
	return err != nil
}

func (w *writeCounter) Try64(n int64, err error) bool {
	w.n += n
	w.err = err
	return err != nil
}

func (w *writeCounter) Result() (int64, error) {
	return w.n, w.err
}

func dumpPair(key string, list []string, w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(io.WriteString(w, key)) || wc.Try(io.WriteString(w, " /")) {
		return wc.Result()
	}
	for _, candidate := range list {
		if wc.Try(io.WriteString(w, candidate)) || wc.Try(io.WriteString(w, "/")) {
			return wc.Result()
		}
	}
	wc.Try(io.WriteString(w, "\n"))
	return wc.Result()
}

// WriteTo outputs the contents of dictonary with UTF8
func (j Jisyo) WriteTo(w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(io.WriteString(w, ";; okuri-ari entries.\n")) {
		return wc.Result()
	}
	for key, list := range j {
		if r, _ := utf8.DecodeLastRuneInString(key); 'a' <= r && r <= 'z' {
			if wc.Try64(dumpPair(key, list, w)) {
				return wc.Result()
			}
		}
	}
	if wc.Try(io.WriteString(w, "\n;; okuri-nasi entries.\n")) {
		return wc.Result()
	}
	for key, list := range j {
		if r, _ := utf8.DecodeLastRuneInString(key); r < 'a' || 'z' < r {
			if wc.Try64(dumpPair(key, list, w)) {
				return wc.Result()
			}
		}
	}
	return wc.Result()
}

// WriteTo outputs the contents of dictonary with EUC-JP
func (j Jisyo) WriteToEucJp(w io.Writer) (n int64, err error) {
	encoder := japanese.EUCJP.NewEncoder()
	return j.WriteTo(encoder.Writer(w))
}
