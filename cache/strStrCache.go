package cache

import (
	"fmt"
	"strings"
)

type StrStrCache map[string]*StringNode

func (c StrStrCache) Update(key string, value string) {
	c[key] = &StringNode{value: value, next: c[key]}
}

func (c StrStrCache) String() string {
	var b strings.Builder
	for key, first := range c {
		b.WriteString(fmt.Sprintf("%s: [%s]\n", key, first.MkString(", ")))
	}
	return b.String()
}

func MakeStrStr() StrStrCache {
	return make(map[string]*StringNode)
}

func (c StrStrCache) Duplicates() StrStrCache {
	var filtered StrStrCache = MakeStrStr()
	for key, first := range c {
		if first.HasNext() {
			filtered[key] = first
		}
	}
	return filtered
}
