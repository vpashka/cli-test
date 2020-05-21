package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

type FakeWriter struct {
	buf []byte
}

func (fw *FakeWriter) Write(data []byte) (int, error) {

	fw.buf = append(fw.buf, data...)
	return len(data), nil
}

func NewFakeWriter() *FakeWriter {
	return &FakeWriter{
		buf: make([]byte, 0, 1024),
	}
}

func check(sourceFile string, resultFile string) error {

	srcData, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("read '%s' error: %s", sourceFile, err)
	}

	result, err := ioutil.ReadFile(resultFile)
	if err != nil {
		return fmt.Errorf("read '%s' error: %s", resultFile, err)
	}

	src, err := parseSource(srcData)
	if err != nil {
		return fmt.Errorf("(main1): parse source error: %s", err)
	}

	nodes, err := buildDocument(src)
	if err != nil {
		return fmt.Errorf("(main1): build document error: %s", err)
	}

	fw := NewFakeWriter()

	err = writeOut(nodes, fw)
	if err != nil {
		return fmt.Errorf("(main1): write error: %s", err)
	}

	if bytes.Compare(result, fw.buf) != 0 {
		return fmt.Errorf("(main1): Output result error:\nWant:\n%s \n\nGot:\n%s", result, fw.buf)
	}

	return nil
}

func Test_main(t *testing.T) {

	err := check("test-data/test1.source.yaml", "test-data/test1.result.yaml")
	if err != nil {
		t.Errorf("(test1): %s", err)
	}

	err = check("test-data/test2.source.yaml", "test-data/test2.result.yaml")
	if err != nil {
		t.Errorf("(test2): %s", err)
	}
}
