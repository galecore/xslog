package internal

import (
	"fmt"
	"testing"

	"github.com/jba/slog/withsupport"
	"golang.org/x/exp/slog"
)

func FuzzWithSupport(f *testing.F) {
	myGroup := NewHelper(false)
	jbaGroup := &withsupport.GroupOrAttrs{}
	iteration := 0
	f.Fuzz(func(t *testing.T, coinFlip int, attrCount int, groupName string) {
		iteration++
		if coinFlip%2 == 0 {
			myGroup = myGroup.WithGroup(groupName)
			jbaGroup = jbaGroup.WithGroup(groupName)
		} else {
			var attrs []slog.Attr
			for i := 0; i < attrCount; i++ {
				attrs = append(attrs, slog.Int(fmt.Sprintf("iter-%d-attr-%d", iteration, i), i))
			}
			myGroup = myGroup.WithAttrs(attrs)
			jbaGroup = jbaGroup.WithAttrs(attrs)
		}

		testEqualGroups(t, myGroup, jbaGroup)
	})
}

type attrValue struct {
	groups []string
	attr   slog.Attr
}

func (a attrValue) Equal(t *testing.T, other attrValue) bool {
	if len(a.groups) != len(other.groups) {
		t.Errorf("groups length mismatch: %v != %v", a.groups, other.groups)
		return false
	}
	for i := range a.groups {
		if a.groups[i] != other.groups[i] {
			t.Errorf("groups mismatch: %v != %v", a.groups, other.groups)
			return false
		}
	}
	if !a.attr.Equal(other.attr) {
		t.Errorf("attr mismatch: %v != %v", a.attr, other.attr)
		return false
	}
	return true
}

func testEqualGroups(t *testing.T, myGroup *Helper, jbaGroup *withsupport.GroupOrAttrs) {
	myGroupAttrs := make(map[string]attrValue)
	myGroup.Apply(func(groups []string, attr slog.Attr) {
		myGroupAttrs[attr.Key] = attrValue{groups, attr}
	})

	jbaGroupAttrs := make(map[string]attrValue)
	jbaGroup.Apply(func(groups []string, attr slog.Attr) {
		jbaGroupAttrs[attr.Key] = attrValue{groups, attr}
	})

	for key, attr := range myGroupAttrs {
		jbaAttr, ok := jbaGroupAttrs[key]
		if !ok {
			t.Errorf("extra key %q", key)
			continue
		}
		if !attr.Equal(t, jbaAttr) {
			t.Errorf("mismatch for key %q: %v != %v", key, attr, jbaAttr)
		}
	}

	for key, attr := range jbaGroupAttrs {
		myAttr, ok := myGroupAttrs[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if !attr.Equal(t, myAttr) {
			t.Errorf("mismatch for key %q: %v != %v", key, attr, myAttr)
		}
	}
}
