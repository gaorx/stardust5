package sderr

import (
	"fmt"
	"strings"
)

type Attr struct {
	Key   string
	Value any
}

type Attrs map[string]any

func NewWith(msg string, attrs ...any) error {
	return New(formatAttrs(msg, attrs))
}

func WrapWith(err error, msg string, attrs ...any) error {
	return Wrap(err, formatAttrs(msg, attrs))
}

func formatAttrs(msg string, attrs []any) string {
	if len(attrs) <= 0 {
		return msg
	}

	var buf strings.Builder
	buf.WriteString(msg)
	first := true
	writeAttr := func(attr *Attr) {
		if first {
			first = false
			buf.WriteString(" [")
		} else {
			buf.WriteString(" ")
		}
		if attr.Key != "" {
			buf.WriteString(attr.Key)
			buf.WriteString("=")
		}
		if attr.Value != nil {
			buf.WriteString(fmt.Sprintf("%v", attr.Value))
		} else {
			buf.WriteString("nil")
		}
	}

	writeAttrs := func(attrs Attrs) {
		for k, v := range attrs {
			writeAttr(&Attr{Key: k, Value: v})
		}
	}

	for _, attr := range attrs {
		switch x := attr.(type) {
		case nil:
			writeAttr(&Attr{Value: nil})
		case Attr:
			writeAttr(&x)
		case *Attr:
			if x != nil {
				writeAttr(x)
			}
		case Attrs:
			writeAttrs(x)
		case map[string]any:
			writeAttrs(x)
		default:
			writeAttr(&Attr{Value: x})
		}
	}
	if !first {
		buf.WriteString("]")
	}
	return buf.String()
}
