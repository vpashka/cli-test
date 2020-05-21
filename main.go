package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// ----------------------------------------------------------------------------
type Node map[string]interface{}

type SourceNode struct {
	Name   string
	Matrix []interface{}
}

// asList - checks if item is a list
func asList(m interface{}, matrix *[]interface{}) bool {

	list, ok := m.([]interface{})
	if ok {
		for _, i := range list {
			*matrix = append(*matrix, i)
		}
	}

	return ok
}

// asMap - checks if item is a map
func asMap(m interface{}, matrix *[]interface{}) bool {

	tmp, ok := m.(map[interface{}]interface{})
	if ok {
		for k, v := range tmp {
			i := make(map[interface{}]interface{})
			i[k] = v
			*matrix = append(*matrix, i)
		}
	}

	return ok
}

// asItem - checks if item is a simple item
func asItem(m interface{}, matrix *[]interface{}) bool {

	i, ok := m.(interface{})
	if ok {
		*matrix = append(*matrix, i)
	}

	return ok
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
		logger.Fatal("Unknown input file. Use: cli input.yaml [output.yaml]")
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

	srcyaml := make(map[string]interface{}, 0)

	err = yaml.Unmarshal([]byte(data), &srcyaml)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	if len(srcyaml) == 0 {
		logger.Fatalf("Not found nodes in '%s'", infile)
	}

	log := logger.WithField("inputfile", infile)

	srcnodes := make([]*SourceNode, 0)

	for name, v := range srcyaml {

		snode := &SourceNode{
			Name:   name,
			Matrix: make([]interface{}, 0),
		}

		srcnodes = append(srcnodes, snode)

		tmp, ok := v.(map[interface{}]interface{})
		if !ok {
			log.Fatalf("can't parse node '%s' as map", name)
		}

		m, ok := tmp["matrix"]
		if !ok {
			log.Fatalf("Not found parameter 'matrix' in node '%s'", name)
		}

		ok = asList(m, &snode.Matrix)
		if ok {
			continue
		}

		ok = asMap(m, &snode.Matrix)
		if ok {
			continue
		}

		ok = asItem(m, &snode.Matrix)
		if ok {
			continue
		}

		log.Fatalf("can't parse matrix for '%s'. Must be list of 'key: value' or 'key: value'", name)
	}

	nodes := make([]Node, 0, len(srcnodes))

	// building the final yaml-document (as map)
	for _, src := range srcnodes {
		for _, item := range src.Matrix {
			node := make(Node, 0)
			node[src.Name] = item
			nodes = append(nodes, node)
		}
	}

	fd := os.Stdout

	if outfile != "" {
		fd, err = os.Create(outfile)
		if err != nil {
			log.Fatalf("Create output file error: %s", err)
		}

		defer fd.Close()
	}

	// writing the final document (as yaml)
	for _, node := range nodes {
		s, err := yaml.Marshal(&node)
		if err != nil {
			log.Fatalf("marshaling error: %s", err)
		}

		s = append(s, '\n')

		_, err = fd.Write(s)
		if err != nil {
			log.Fatalf("write output error: %s", err)
		}
	}
}
