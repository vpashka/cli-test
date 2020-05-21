package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// writeOut - write final yaml document
func writeOut(nodes []Node, out io.Writer) error {

	for k, node := range nodes {
		s, err := yaml.Marshal(&node)
		if err != nil {
			return fmt.Errorf("marshaling error: %s", err)
		}

		// line feed between nodes
		if k > 0 {
			_, err = out.Write([]byte{'\n'})
			if err != nil {
				return fmt.Errorf("write output error: %s", err)
			}
		}

		// write data
		n, err := out.Write(s)
		if err != nil {
			return fmt.Errorf("write output error: %s", err)
		}

		// TODO: may be repeat write?
		if n < len(s) {
			return fmt.Errorf("write data error: want=%d bytes real write=%d bytes", len(s), n)
		}
	}

	return nil
}
