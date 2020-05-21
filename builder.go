package main

// ----------------------------------------------------------------------------

// buildDocument - 	building the final yaml-document (as map)
func buildDocument(source []*SourceNode) ([]Node, error) {

	nodes := make([]Node, 0, len(source))

	for _, src := range source {
		for _, item := range src.Matrix {
			node := newNode()
			node[src.Name] = item
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}
