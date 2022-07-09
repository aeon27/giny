package giny

import (
	"reflect"
	"testing"
)

func testParsePattern(s string, expected []string, t *testing.T) bool {
	result := parsePattern(s)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("test parsePattern failed: %s", s)
		return false
	}
	return true
}

func TestParsePattern(t *testing.T) {
	testParsePattern("/p/:name", []string{"p", ":name"}, t)
	testParsePattern("/p/*", []string{"p", "*"}, t)
	testParsePattern("/p/*name/*", []string{"p", "*name"}, t)
}
