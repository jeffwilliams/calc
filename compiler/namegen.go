package compiler

import "fmt"

type nameGenerator struct {
	prefix string
	next   int
}

func newNameGenerator(prefix string) nameGenerator {
	return nameGenerator{prefix: prefix}
}

func (n *nameGenerator) alloc() string {
	s := fmt.Sprintf(n.prefix, n.next)
	n.next++
	return s

}
