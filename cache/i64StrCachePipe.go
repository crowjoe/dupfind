package cache

import (
	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/file"
)

func (c *I64StrCache) ToChan(ctl *channel.Control) chan *file.File {
	out := make(chan *file.File)
	go c.pipeTo(out, ctl)
	return out
}

func (c *I64StrCache) pipeTo(out chan<- *file.File, ctl *channel.Control) {
	defer close(out)

	for size, nodes := range *c {
		var node *StringNode = nodes
		for node != nil {
			if ctl.Stopped() {
				return
			}
			f := file.File{Name: node.value, Size: size}
			out <- &f
			node = node.next
		}
	}
}
