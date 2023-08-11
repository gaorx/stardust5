package sdreq

import (
	"fmt"

	"github.com/samber/lo"
)

func ToQueryParam(v any) string {
	if v == nil {
		return ""
	}
	if v1, ok := v.(string); ok {
		return v1
	}
	if v1, ok := v.(fmt.Stringer); ok {
		return v1.String()
	}
	return fmt.Sprintf("%v", v)
}

func ToQueryParams(params map[string]any) map[string]string {
	return lo.MapValues(params, func(v any, _ string) string {
		return ToQueryParam(v)
	})
}
