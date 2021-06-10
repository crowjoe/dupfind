package cache

// NOTE - to run all tests under current dir: `go test ./...`

import (
	"testing"
)

func TestNodeHasNext(t *testing.T) {
	n1 := &StringNode{value: "foo"}
	if n1.HasNext() {
		t.Errorf("HasNext should be false, but was true")
	}

	n2 := &StringNode{value: "bar", next: n1}
	if !n2.HasNext() {
		t.Errorf("HasNext should be true, but was false")
	}

}

func TestNodeMkString(t *testing.T) {
	n1 := &StringNode{value: "foo"}
	n2 := &StringNode{value: "bar", next: n1}

	s1 := n1.MkString(", ")
	e1 := "foo"
	if s1 != e1 {
		t.Errorf("s1 had value %s instead of %s", s1, e1)
	}

	s2 := n2.MkString(", ")
	e2 := "bar, foo"
	if s2 != e2 {
		t.Errorf("s2 had value %s instead of %s", s2, e2)
	}
}
