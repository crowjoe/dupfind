package cache

import (
	"strings"
	"testing"
)

func TestI64StrCacheUpdate(t *testing.T) {
	var cache I64StrCache = MakeI64Str()

	cache.Update(1, "value1")

	n1 := cache[1]

	if n1 == nil {
		t.Errorf("nothing found for checksum1")
	}
	if n1 != nil && n1.next != nil {
		t.Errorf("next should be nil")
	}
	if n1 != nil && n1.value != "value1" {
		t.Errorf("value should be 'value1' but was %s", n1.value)
	}

	cache.Update(1, "value2")

	n1 = cache[1]

	if n1.value != "value2" {
		t.Errorf("n1.value should be value2")
	}
	if n1.next == nil {
		t.Errorf("next should NOT be nil")
	}
	if n1.next.value != "value1" {
		t.Errorf("second value should be 'value1' but was %s", n1.next.value)
	}

	cache.Update(2, "file3")

	n2 := cache[2]

	if n2 == nil {
		t.Errorf("nothing found for checksum2")
	}
	if n2 != nil && n2.next != nil {
		t.Errorf("next should be nil")
	}
	if n2 != nil && n2.value != "file3" {
		t.Errorf("value should be 'file3' but was %s", n2.value)
	}

}

func TestI64StrCacheString(t *testing.T) {
	var cache I64StrCache = MakeI64Str()

	cache.Update(1, "value1")
	cache.Update(1, "value2")
	cache.Update(2, "file3")

	results := cache.String()

	expected1 := "1: [value2, value1]\n"
	expected2 := "2: [file3]\n"

	if !strings.Contains(results, expected1) {
		t.Errorf("'%s' did not include '%s'", results, expected1)
	}

	if !strings.Contains(results, expected2) {
		t.Errorf("'%s' did not include '%s'", results, expected2)
	}
}

func TestI64StrCacheDuplicates(t *testing.T) {
	var cache I64StrCache = MakeI64Str()

	cache.Update(1, "value1")
	cache.Update(1, "value2")
	cache.Update(2, "file3")

	if cache[2] == nil {
		t.Errorf("expected non-nil for checksum2")
	}

	var duplicates = cache.Duplicates()

	if duplicates[2] != nil {
		t.Errorf("expected nil for checksum2")
	}

	if cache[1] == nil {
		t.Errorf("expected non-nil for checksum1")
	}

}

func TestI64StrCacheDuplicatesString(t *testing.T) {
	var cache I64StrCache = MakeI64Str()

	cache.Update(1, "value1")
	cache.Update(1, "value2")
	cache.Update(2, "file3")

	var duplicates = cache.Duplicates()

	results := duplicates.String()

	var expected string = "1: [value2, value1]\n"

	if results != expected {
		t.Errorf("expected '%s' got '%s'", expected, results)
	}
}
