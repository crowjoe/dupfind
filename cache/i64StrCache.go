package cache

import (
	"fmt"
	"strings"
)

type I64StrCache map[int64]*StringNode

func (c I64StrCache) Update(key int64, value string) {
	c[key] = &StringNode{value: value, next: c[key]}
}

func (c I64StrCache) String() string {
	var b strings.Builder
	for key, first := range c {
		b.WriteString(fmt.Sprintf("%d: [%s]\n", key, first.MkString(", ")))
	}
	return b.String()
}

func MakeI64Str() I64StrCache {
	return make(map[int64]*StringNode)
}

func (c I64StrCache) Duplicates() I64StrCache {
	var filtered I64StrCache = MakeI64Str()
	for key, first := range c {
		if first.HasNext() {
			filtered[key] = first
		}
	}
	return filtered
}
