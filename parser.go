package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// ----------------------------------------------------------------------------
// asList - checks if item is a list
func asList(m interface{}, snode *SourceNode) bool {

	list, ok := m.([]interface{})
	if ok {
		for _, i := range list {
			snode.Matrix = append(snode.Matrix, i)
		}
	}

	return ok
}

// asMap - checks if item is a map
func asMap(m interface{}, snode *SourceNode) bool {

	tmp, ok := m.(map[interface{}]interface{})
	if ok {
		for k, v := range tmp {
			i := make(map[interface{}]interface{})
			i[k] = v
			snode.Matrix = append(snode.Matrix, i)
		}
	}

	return ok
}

// asItem - checks if item is a simple item
func asItem(m interface{}, snode *SourceNode) bool {

	i, ok := m.(interface{})
	if ok {
		snode.Matrix = append(snode.Matrix, i)
	}

	return ok
}

// ----------------------------------------------------------------------------
// parseSource - parse data as a yaml document
func parseSource(data []byte) ([]*SourceNode, error) {

	yamldoc := make(map[string]interface{}, 0)

	err := yaml.Unmarshal(data, &yamldoc)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err)
	}

	if len(yamldoc) == 0 {
		return nil, fmt.Errorf("Not found nodes")
	}

	srcnodes := make([]*SourceNode, 0)

	for name, v := range yamldoc {

		snode := &SourceNode{
			Name:   name,
			Matrix: make([]interface{}, 0),
		}

		srcnodes = append(srcnodes, snode)

		tmp, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("can't parse node '%s' as map", name)
		}

		m, ok := tmp["matrix"]
		if !ok {
			return nil, fmt.Errorf("Not found parameter 'matrix' in node '%s'", name)
		}

		ok = asList(m, snode)
		if ok {
			continue
		}

		ok = asMap(m, snode)
		if ok {
			continue
		}

		ok = asItem(m, snode)
		if ok {
			continue
		}

		return nil, fmt.Errorf("can't parse matrix for '%s'. Must be list of 'key: value' or 'key: value'", name)
	}

	return srcnodes, nil
}
