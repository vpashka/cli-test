package main

import (
	"testing"
)

func Test_parseSource(t *testing.T) {

	data1 := []byte(`
foo:
  matrix:
    name1: bar
    name2: baz
`)

	src, err := parseSource(data1)
	if err != nil {
		t.Errorf("(parseSource): test1: parse error: %s", err)
		return
	}

	if len(src) != 1 {
		t.Errorf("(parseSource): test1: size of source error: want=1 got=%d", len(src))
		return
	}

	if len(src[0].Matrix) != 2 {
		t.Errorf("(parseSource): test1: 'matrix' size error: want=2 got=%d", len(src[0].Matrix))
		return
	}

	data2 := []byte(`
node1:
  matrix:
   - config:
       foo: bar
   - key: value
   - key2: value2
`)

	src, err = parseSource(data2)
	if err != nil {
		t.Errorf("(parseSource): test2: parse error: %s", err)
		return
	}

	if len(src) != 1 {
		t.Errorf("(parseSource): test2: size of source error: want=1 got=%d", len(src))
		return
	}

	if len(src[0].Matrix) != 3 {
		t.Errorf("(parseSource): test2: 'matrix' size error: want=3 got=%d", len(src[0].Matrix))
		return
	}

	data3 := []byte(`
foo:
  matrix:
    - name: bar
    - name: baz

foo2:
  matrix:
    - bazname: bar
    - bazname: baz   
`)

	src, err = parseSource(data3)
	if err != nil {
		t.Errorf("(parseSource): test3: parse error: %s", err)
		return
	}

	if len(src) != 2 {
		t.Errorf("(parseSource): test3: size of source error: want=2 got=%d", len(src))
		return
	}

	if len(src[1].Matrix) != 2 {
		t.Errorf("(parseSource): test3: 'matrix' size error: want=2 got=%d", len(src[1].Matrix))
		return
	}
}

func Test_badSource(t *testing.T) {

	data1 := []byte(`
foo:
`)

	_, err := parseSource(data1)
	if err == nil {
		t.Errorf("(parseBadSource): parse 'bad structure' error: want=error got=nil")
		return
	}

	data2 := []byte(``)

	_, err = parseSource(data2)
	if err == nil {
		t.Errorf("(parseBadSource): parse 'empty structure' error: want=error got=nil")
		return
	}

	data3 := []byte(`
foo:
  matrix:
`)

	_, err = parseSource(data3)
	if err == nil {
		t.Errorf("(parseBadSource): parse 'empty matrix' error: want=error got=nil")
		return
	}
}
