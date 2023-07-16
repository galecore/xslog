package internal

import (
	"golang.org/x/exp/slog"
)

type Helper struct {
	parent *Helper

	groupCache []string

	name  string
	attrs []slog.Attr

	flattenAttrGroups bool
}

func NewHelper(flattenAttrGroups bool) *Helper {
	return &Helper{flattenAttrGroups: flattenAttrGroups}
}

func (h *Helper) WithAttrs(attrs []slog.Attr) *Helper {
	return &Helper{
		parent:            h,
		attrs:             attrs,
		flattenAttrGroups: h.flattenAttrGroups,
	}
}

func (h *Helper) WithGroup(name string) *Helper {
	return &Helper{
		parent:            h,
		name:              name,
		flattenAttrGroups: h.flattenAttrGroups,
	}
}

func (h *Helper) groups() []string {
	if h.parent == nil {
		if h.name != "" {
			h.groupCache = []string{h.name}
		}
		return h.groupCache
	}
	if h.groupCache != nil {
		return h.groupCache
	}
	h.groupCache = h.parent.groups()
	if h.name != "" {
		h.groupCache = append(h.groupCache, h.name)
	}
	return h.groupCache
}

func (h *Helper) Apply(f func(groups []string, attr slog.Attr)) {
	groups := h.groups()
	for p := h.parent; p != nil; p = p.parent {
		p.applyToNode(f, p.groupCache, p.attrs)
	}
	h.applyToNode(f, groups, h.attrs)
}

func (h *Helper) applyToNode(f func(groups []string, attr slog.Attr), groups []string, attrs []slog.Attr) {
	for _, attr := range attrs {
		if attr.Value.Kind() == slog.KindGroup && h.flattenAttrGroups {
			h.applyToNode(f, append(groups, attr.Key), attr.Value.Group())
		} else {
			f(groups, attr)
		}
	}
}
