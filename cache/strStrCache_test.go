package cache

import (
	"strings"
	"testing"
)

func TestStrStrCacheUpdate(t *testing.T) {
	var cache StrStrCache = MakeStrStr()

	cache.Update("key1", "value1")

	n1 := cache["key1"]

	if n1 == nil {
		t.Errorf("nothing found for checksum1")
	}
	if n1 != nil && n1.next != nil {
		t.Errorf("next should be nil")
	}
	if n1 != nil && n1.value != "value1" {
		t.Errorf("value should be 'file1' but was %s", n1.value)
	}

	cache.Update("key1", "value2")

	n1 = cache["key1"]

	if n1.value != "value2" {
		t.Errorf("n1.value should be file2")
	}
	if n1.next == nil {
		t.Errorf("next should NOT be nil")
	}
	if n1.next.value != "value1" {
		t.Errorf("second value should be 'file1' but was %s", n1.next.value)
	}

	cache.Update("key2", "value3")

	n2 := cache["key2"]

	if n2 == nil {
		t.Errorf("nothing found for checksum2")
	}
	if n2 != nil && n2.next != nil {
		t.Errorf("next should be nil")
	}
	if n2 != nil && n2.value != "value3" {
		t.Errorf("value should be 'file3' but was %s", n2.value)
	}

}

func TestStrStrCacheString(t *testing.T) {
	var cache StrStrCache = MakeStrStr()

	cache.Update("key1", "value1")
	cache.Update("key1", "value2")
	cache.Update("key2", "value3")

	results := cache.String()

	expected1 := "key1: [value2, value1]\n"
	expected2 := "key2: [value3]\n"

	if !strings.Contains(results, expected1) {
		t.Errorf("'%s' did not include '%s'", results, expected1)
	}

	if !strings.Contains(results, expected2) {
		t.Errorf("'%s' did not include '%s'", results, expected2)
	}
}

func TestStrStrCacheDuplicates(t *testing.T) {
	var cache StrStrCache = MakeStrStr()

	cache.Update("key1", "value1")
	cache.Update("key1", "value2")
	cache.Update("key2", "value3")

	if cache["key2"] == nil {
		t.Errorf("expected non-nil for checksum2")
	}

	var duplicates = cache.Duplicates()

	if duplicates["key2"] != nil {
		t.Errorf("expected nil for checksum2")
	}

	if cache["key1"] == nil {
		t.Errorf("expected non-nil for checksum1")
	}

}

func TestStrStrCacheDuplicatesString(t *testing.T) {
	var cache StrStrCache = MakeStrStr()

	cache.Update("key1", "value1")
	cache.Update("key1", "value2")
	cache.Update("key2", "value3")

	var duplicates = cache.Duplicates()

	results := duplicates.String()

	var expected string = "key1: [value2, value1]\n"

	if results != expected {
		t.Errorf("expected '%s' got '%s'", expected, results)
	}
}
