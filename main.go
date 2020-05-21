package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// ----------------------------------------------------------------------------
type Node map[string]interface{}

func newNode() Node {
	return make(map[string]interface{}, 0)
}

type SourceNode struct {
	Name   string
	Matrix []interface{}
}

// ----------------------------------------------------------------------------
func main() {

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&nested.Formatter{
		HideKeys: true,
		// DisableColors: true,
	})

	if len(os.Args) < 2 {
		logger.Fatalf("Unknown input file. Use: %s input.yaml [output.yaml]", os.Args[0])
	}

	infile := os.Args[1]

	outfile := ""
	if len(os.Args) >= 3 {
		outfile = os.Args[2]
	}

	data, err := ioutil.ReadFile(infile)
	if err != nil {
		logger.Fatalf("open '%s' error: %s", infile, err)
	}

	log := logger.WithField("inputfile", infile)

	srcnodes, err := parseSource(data)
	if err != nil {
		log.Fatalf("parse source error: %s", err)
	}

	nodes, err := buildDocument(srcnodes)
	if err != nil {
		log.Fatalf("build yaml error: %s", err)
	}

	// writing the final document (as yaml)
	if outfile == "" {
		err = writeOut(nodes, os.Stdout)
		if err != nil {
			log.Errorf("%s", err)
		}
	} else {

		fd, err := os.Create(outfile)
		if err != nil {
			log.Fatalf("Create output file error: %s", err)
		}
		defer fd.Close()

		err = writeOut(nodes, fd)
		if err != nil {
			log.Errorf("%s", err)
		}
	}
}
