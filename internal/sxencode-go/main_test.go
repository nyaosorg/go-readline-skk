package sxencode

import (
	"reflect"
	"strings"
	"testing"
)

type foo struct {
	Bar   string
	Corge func()
}

func TestStruct(t *testing.T) {
	v := &foo{
		Bar:   "hoge",
		Corge: func() {},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `((Bar "hoge"))`
	result := string(b)
	if expect != result {
		t.Fatalf("expect %v, but %v", expect, result)
	}

	var sbuf strings.Builder
	enc := NewEncoder(&sbuf)
	enc.OnTypeNotSupported = func(v reflect.Value) (string, error) {
		return "not-support-type", nil
	}
	enc.Encode(v)
	result = sbuf.String()
	expect1 := `((Bar "hoge")(Corge not-support-type))`
	expect2 := `((Corge not-support-type)(Bar "hoge"))`
	if expect1 != result && expect2 != result {
		t.Fatalf("expect %v or %v, but %v", expect1, expect2, result)
	}
}

func TestMap(t *testing.T) {
	v := map[string]any{
		"bar": "hoge",
		"baz": func() {},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `(("bar" "hoge"))`
	result := string(b)
	if expect != result {
		t.Fatalf("expect %v, but %v", expect, result)
	}

	var sbuf strings.Builder
	enc := NewEncoder(&sbuf)
	enc.OnTypeNotSupported = func(v reflect.Value) (string, error) {
		return "not-support-type", nil
	}
	enc.Encode(v)
	result = sbuf.String()
	expect1 := `(("bar" "hoge")("baz" not-support-type))`
	expect2 := `(("baz" not-support-type)("bar" "hoge"))`
	if expect1 != result && expect2 != result {
		t.Fatalf("expect %v or %v, but %v", expect1, expect2, result)
	}
}

func TestStructWithTag(t *testing.T) {
	type fooWithTag struct {
		Bar string `sxpr:"bar-field"`
	}
	v := &fooWithTag{
		Bar: "value",
	}
	s, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `((bar-field "value"))`
	result := string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}
}

func TestStructOmit(t *testing.T) {
	type fooWithTag struct {
		Bar string `sxpr:"bar,omitempty"`
	}
	v := &fooWithTag{}
	s, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `()`
	result := string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}

	v = &fooWithTag{Bar: "x"}
	s, err = Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect = `((bar "x"))`
	result = string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}
}

func TestStructOmitOnly(t *testing.T) {
	type fooWithTag struct {
		Bar string `sxpr:",omitempty"`
		Baz string `sxpr:"baz"`
	}
	v := &fooWithTag{Baz: "1"}
	s, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `((baz "1"))`
	result := string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}

	v = &fooWithTag{Bar: "x", Baz: "1"}
	s, err = Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect1 := `((Bar "x")(baz "1"))`
	expect2 := `((baz "1")(Bar "x"))`
	result = string(s)
	if expect1 != result && expect2 != result {
		t.Fatalf("expect %#v or %#v, but %#v", expect1, expect2, result)
	}
}

func TestHyphen(t *testing.T) {
	type foo struct {
		Bar string `sxpr:"bar"`
		Baz string `sxpr:"-"`
	}
	v := &foo{
		Bar: "1",
		Baz: "2",
	}
	s, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `((bar "1"))`
	result := string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}
}

func TestNoName(t *testing.T) {
	type foo struct {
		Bar string `sxpr:",noname"`
		Baz int    `sxpr:"baz"`
	}
	v := &foo{
		Bar: "one",
		Baz: 7,
	}
	s, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `( "one"(baz 7))`
	result := string(s)
	if expect != result {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}
}
