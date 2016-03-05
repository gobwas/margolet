package slugger

import (
	"reflect"
	"testing"
)

func TestSluggerFill(t *testing.T) {
	var s []string
	fill(&s, "a b")
	if !reflect.DeepEqual(s, []string{"a", "b"}) {
		t.Errorf("fill([]string, \"a b\")=%v; want %v", s, []string{"a", "b"})
	}
}
