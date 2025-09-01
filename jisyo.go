package skk

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
)

const (
	ariHeader  = ";; okuri-ari entries."
	nasiHeader = ";; okuri-nasi entries."
)

type _History struct {
	key string
	val []candidateT
}

type candidateT interface {
	String() string
	Source() string
}

type candidateStringT string

func (c candidateStringT) String() string { return string(c) }

var encodeCandidate = strings.NewReplacer(
	`/`, `\057`,
	`"`, `\"`,
	`\`, `\\`,
)

func (c candidateStringT) Source() string {
	s := string(c)
	if strings.ContainsRune(s, '/') {
		return fmt.Sprintf(`(concat "%s")`, encodeCandidate.Replace(s))
	}
	return s
}

type candidateFuncT struct {
	source string
	f      func() string
}

func (c *candidateFuncT) Source() string { return c.source }
func (c *candidateFuncT) String() string { return c.f() }

// Jisyo is a dictionary that contains user or system dictionary.
type Jisyo struct {
	ari         map[string][]candidateT
	nasi        map[string][]candidateT
	ariHistory  []_History
	nasiHistory []_History
}

func newJisyo() *Jisyo {
	return &Jisyo{
		ari:  map[string][]candidateT{},
		nasi: map[string][]candidateT{},
	}
}

func (j *Jisyo) lookup(key string, okuri bool) (candidates []candidateT, ok bool) {
	if okuri {
		candidates, ok = j.ari[key]
	} else {
		candidates, ok = j.nasi[key]
	}
	return
}

func (j *Jisyo) store(key string, okuri bool, value []candidateT) {
	if okuri {
		j.ari[key] = value
	} else {
		j.nasi[key] = value
	}
}

func (j *Jisyo) storeAndLearn(key string, okuri bool, value []candidateT) {
	j.store(key, okuri, value)
	if okuri {
		j.ariHistory = append(j.ariHistory, _History{key: key, val: value})
	} else {
		j.nasiHistory = append(j.nasiHistory, _History{key: key, val: value})
	}
}

func (j *Jisyo) remove(key string, okuri bool) {
	if okuri {
		delete(j.ari, key)
		j.ariHistory = append(j.ariHistory, _History{key: key, val: nil})
	} else {
		delete(j.nasi, key)
		j.nasiHistory = append(j.nasiHistory, _History{key: key, val: nil})
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
// It returns time-stamp and error.
func (j *Jisyo) Load(filename string) error {
	filename = expandEnv(filename)
	matches, err := filepath.Glob(filename)
	if err != nil {
		_, err = j.load(filename)
		return err
	}
	for _, fn := range matches {
		if _, err = j.load(fn); err != nil {
			return err
		}
	}
	return nil
}

func (j *Jisyo) load(filename string) (time.Time, error) {
	var stamp time.Time
	fd, err := os.Open(filename)
	if err != nil {
		return stamp, err
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err == nil {
		stamp = stat.ModTime()
	}
	return stamp, j.Read(fd)
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
			if len(one) > 2 && one[0] == '(' && one[len(one)-1] == ')' {
				values = append(values, evalSxString(one))
			} else {
				values = append(values, candidateStringT(one))
			}
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

func peekLine(r io.Reader) (io.Reader, string, error) {
	var buffer bytes.Buffer
	br := bufio.NewReader(r)

	line, err := br.ReadString('\n')
	buffer.WriteString(line)
	io.CopyN(&buffer, br, int64(br.Buffered()))
	if err == io.EOF {
		return &buffer, line, nil
	}
	return io.MultiReader(&buffer, r), line, err
}

func (j *Jisyo) Read(r io.Reader) error {
	utf8mode := false
	r, line, err := peekLine(r)
	if err != nil {
		return err
	}
	if len(line) > 0 && line[0] == ';' {
		if m := pragma(line[1:]); m != nil && m["coding"] == "utf-8" {
			utf8mode = true
		}
	}
	if !utf8mode {
		r = japanese.EUCJP.NewDecoder().Reader(r)
	}
	sc := bufio.NewScanner(r)
	okuri := false
	for sc.Scan() {
		okuri = j.readOne(sc.Text(), okuri)
	}
	return sc.Err()
}

func dumpPair(key string, list []candidateT, w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(io.WriteString(w, key)) || wc.Try(io.WriteString(w, " /")) {
		return wc.Result()
	}
	for _, candidate := range list {
		if wc.Try(io.WriteString(w, candidate.Source())) {
			return wc.Result()
		}
		if wc.Try(io.WriteString(w, "/")) {
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

func (j *Jisyo) saveAs(fname string) error {
	fd, err := os.Create(fname)
	if err != nil {
		return err
	}
	if _, err := j.writeToUtf8(fd); err != nil {
		fd.Close()
		return err
	}
	return fd.Close()
}

func (j *Jisyo) writeToUtf8(w io.Writer) (n int64, err error) {
	var wc writeCounter
	if wc.Try(fmt.Fprintln(w, ";; -*- mode: fundamental; coding: utf-8 -*-")) {
		return wc.Result()
	}
	wc.Try64(j.writeTo(w))
	return wc.Result()
}
