package cache

import (
	"strings"
)

type StringNode struct {
	value string
	next  *StringNode
}

func (n *StringNode) HasNext() bool {
	return n.next != nil
}

func (n *StringNode) MkString(sep string) string {
	var b strings.Builder
	if n != nil {
		next := n
		b.WriteString(next.value)
		next = n.next
		for next != nil {
			b.WriteString(sep)
			b.WriteString(next.value)
			next = next.next
		}
	}
	return b.String()
}
