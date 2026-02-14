package sxencode

import (
	"testing"
)

func TestDecodeInt(t *testing.T) {
	var v int
	err := Unmarshal([]byte(`1234`), &v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := 1234
	if v != expect {
		t.Fatalf("expect %#v, but %#v", expect, v)
	}
}

func TestDecodeStruct(t *testing.T) {
	type Qux struct {
		Quux string `sxpr:"quux"`
	}
	type Foo struct {
		Bar string `sxpr:"bar"`
		Baz int
		Qux *Qux `sxpr:"qux"`
	}
	var foo Foo

	err := Unmarshal([]byte(`((bar "10")(Baz 4)(qux ((quux "quuux"))))`), &foo)
	if err != nil {
		t.Fatal(err.Error())
	}
	if expect := "10"; foo.Bar != expect {
		t.Fatalf("Bar: expect %#v, but %#v", expect, foo.Bar)
	}
	if expect := 4; foo.Baz != expect {
		t.Fatalf("Baz: expect %#v, but %#v", expect, foo.Baz)
	}
	if expect := "quuux"; foo.Qux.Quux != expect {
		t.Fatalf("Baz: expect %#v, but %#v", expect, foo.Qux.Quux)
	}
}

func TestDecodeMap(t *testing.T) {
	m := map[string]string{}
	err := Unmarshal([]byte(`(("ahaha" "ihihi")("ufufu" "ohoho"))`), m)
	if err != nil {
		t.Fatal(err.Error())
	}
	testcase := [][2]string{
		{"ahaha", "ihihi"},
		{"ufufu", "ohoho"},
	}
	for _, p := range testcase {
		result := m[p[0]]
		expect := p[1]
		if expect != result {
			t.Fatalf("expect %#v, but %#v", expect, result)
		}
	}
}

func TestDecodeStructMap(t *testing.T) {
	type Foo struct {
		Map map[string]int `sxpr:"map"`
	}
	foo := &Foo{}
	err := Unmarshal([]byte(`((map (("a" 1)("b" 2)("c" 3))))`), &foo)
	if err != nil {
		t.Fatal(err.Error())
	}
	testcase := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	for key, val := range testcase {
		result := foo.Map[key]
		expect := val
		if result != expect {
			t.Fatalf("expect %#v, but %#v", expect, result)
		}
	}
}

func TestDecodeSlice(t *testing.T) {
	foo := []int{}
	err := Unmarshal([]byte(`#(9 8 7 6 5 4 3 2 1)`), &foo)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 0; i < 9; i++ {
		expect := 9 - i
		result := foo[i]
		if expect != result {
			t.Fatalf("expect %#v, but %#v", expect, result)
		}
	}

	foo = []int{}
	err = Unmarshal([]byte(`(10 9 8 7 6 5 4 3 2)`), &foo)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 0; i < 9; i++ {
		expect := 10 - i
		result := foo[i]
		if expect != result {
			t.Fatalf("expect %#v, but %#v", expect, result)
		}
	}

}

func TestBoth(t *testing.T) {
	type Baz struct {
		Qux  []int           `sxpr:"qux"`
		Quux map[string]*Baz `sxpr:"quux"`
	}

	type Foo struct {
		Bar string `sxpr:"bar"`
		Baz Baz    `sxpr:"baz"`
	}

	v := &Foo{
		Bar: "bar1",
		Baz: Baz{
			Qux: []int{4, 3, 2, 1},
			Quux: map[string]*Baz{
				"hoge": &Baz{
					Qux: []int{3, 4, 5, 6},
				},
			},
		},
	}
	bin, err := Marshal(&v)
	if err != nil {
		t.Fatal(err.Error())
	}
	// println(string(bin))

	var v2 Foo

	err = Unmarshal(bin, &v2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if v2.Bar != "bar1" {
		t.Fatal("! v2.Bar")
	}
	if v2.Baz.Qux[2] != 2 {
		t.Fatal("! v2.Baz.Qux")
	}
	if v2.Baz.Quux["hoge"].Qux[3] != 6 {
		t.Fatal("! v2.Baz.Quux[hoge].Qux[3]")
	}
}

func TestDecodeNoname(t *testing.T) {
	type foo struct {
		Bar string `sxpr:"bar,noname"`
		Baz string `sxpr:"baz,noname"`
		Qux string `sxpr:"qux"`
	}
	var foo1 foo

	err := Unmarshal([]byte(`("first" "second" (qux "third"))`), &foo1)
	if err != nil {
		t.Fatal(err.Error())
	}
	if expect := "third"; foo1.Qux != expect {
		t.Fatalf("expect %#v, but %#v", expect, foo1.Qux)
	}
	if expect := "first"; foo1.Bar != expect {
		t.Fatalf("expect %#v, but %#v", expect, foo1.Bar)
	}
	if expect := "second"; foo1.Baz != expect {
		t.Fatalf("expect %#v, but %#v", expect, foo1.Baz)
	}
}
