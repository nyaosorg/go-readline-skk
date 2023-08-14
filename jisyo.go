package skk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

const (
	ariHeader  = ";; okuri-ari entries."
	nasiHeader = ";; okuri-nasi entries."
)

// Jisyo is a dictionary that contains user or system dictionary.
type Jisyo struct {
	ari  map[string][]string
	nasi map[string][]string
}

func newJisyo() *Jisyo {
	return &Jisyo{
		ari:  map[string][]string{},
		nasi: map[string][]string{},
	}
}

func (j *Jisyo) lookup(key string, okuri bool) (candidates []string, ok bool) {
	if okuri {
		candidates, ok = j.ari[key]
	} else {
		candidates, ok = j.nasi[key]
	}
	return
}

func (j *Jisyo) store(key string, okuri bool, value []string) {
	if okuri {
		j.ari[key] = value
	} else {
		j.nasi[key] = value
	}
}

func (j *Jisyo) remove(key string, okuri bool) {
	if okuri {
		delete(j.ari, key)
	} else {
		delete(j.nasi, key)
	}
}

var percentEnv = regexp.MustCompile(`%.*?%`)

func expandEnv(s string) string {
	if len(s) > 0 && s[0] == '~' {
		if u, err := user.Current(); err == nil {
			s = u.HomeDir + s[1:]
		}
	}
	return percentEnv.ReplaceAllStringFunc(s, func(m string) string {
		name := m[1 : len(m)-1]
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return m
	})
}

// Load reads the contents of an dictionary from a file.
func (j *Jisyo) Load(filename string) error {
	fd, err := os.Open(expandEnv(filename))
	if err != nil {
		return err
	}
	defer fd.Close()
	return j.Read(fd)
}

func (j *Jisyo) readOne(line string, okuri bool) bool {
	if len(line) > 2 && line[0] == ';' && line[1] == ';' {
		if strings.HasPrefix(line, ariHeader) {
			okuri = true
		} else if strings.HasPrefix(line, nasiHeader) {
			okuri = false
		}
		return okuri
	}
	if len(line) <= 0 || line[0] == ';' {
		return okuri
	}
	source, lists, ok := strings.Cut(line, " /")
	if !ok {
		return okuri
	}
	values, _ := j.lookup(source, okuri)
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
	j.store(source, okuri, values)
	return okuri
}

func pragma(line string) map[string]string {
	_, body, ok := strings.Cut(line, "-*-")
	if !ok {
		return nil
	}
	body, _, ok = strings.Cut(body, "-*-")
	if !ok {
		return nil
	}
	m := map[string]string{}
	for ok {
		var token string
		token, body, ok = strings.Cut(body, ";")
		if key, value, valid := strings.Cut(token, ":"); valid {
			m[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return m
}

func (j *Jisyo) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	decoder := japanese.EUCJP.NewDecoder()
	f := func(s string) string {
		if utf8, err := decoder.String(s); err == nil {
			return utf8
		}
		return s
	}
	var okuri bool
	if sc.Scan() {
		line := f(sc.Text())
		if len(line) > 0 && line[0] == ';' {
			if m := pragma(line[1:]); m != nil && m["coding"] == "utf-8" {
				f = func(s string) string {
					return s
				}
			}
		} else {
			okuri = j.readOne(line, okuri)
		}
	}

	for sc.Scan() {
		text := sc.Text()
		okuri = j.readOne(f(text), okuri)
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
func (j *Jisyo) writeTo(w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(fmt.Fprintln(w, ariHeader)) {
		return wc.Result()
	}
	for key, list := range j.ari {
		if wc.Try64(dumpPair(key, list, w)) {
			return wc.Result()
		}
	}
	if wc.Try(fmt.Fprintf(w, "\n%s\n", nasiHeader)) {
		return wc.Result()
	}
	for key, list := range j.nasi {
		if wc.Try64(dumpPair(key, list, w)) {
			return wc.Result()
		}
	}
	return wc.Result()
}

// WriteTo outputs the contents of dictonary with EUC-JP
func (j *Jisyo) writeToEucJp(w io.Writer) (n int64, err error) {
	var wc writeCounter
	encoder := japanese.EUCJP.NewEncoder()
	if wc.Try(fmt.Fprintln(w, ";; -*- mode: fundamental; coding: euc-jp -*-")) {
		return wc.Result()
	}
	wc.Try64(j.writeTo(encoder.Writer(w)))
	return wc.Result()
}

func (j *Jisyo) writeToUtf8(w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(fmt.Fprintln(w, ";; -*- mode: fundamental; coding: utf-8 -*-")) {
		return wc.Result()
	}
	wc.Try64(j.writeTo(w))
	return wc.Result()
}
