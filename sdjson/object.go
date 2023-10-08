package sdjson

type Object map[string]any

func (o Object) Len() int {
	return len(o)
}

func (o Object) Has(k string) bool {
	_, ok := o[k]
	return ok
}

func (o Object) Get(k string) Value {
	v0, ok := o[k]
	if ok {
		return V(v0)
	} else {
		return V(nil)
	}
}

func (o Object) Set(k string, v any) Object {
	if o != nil {
		o[k] = unbox(v)
	}
	return o
}

func (o Object) TryGet(keys ...string) Value {
	for _, k := range keys {
		v0, ok := o[k]
		if ok {
			return V(v0)
		}
	}
	return V(nil)
}

func (o Object) ShadowCopy() Object {
	if o == nil {
		return nil
	}
	o1 := Object{}
	for k, v := range o {
		o1[k] = v
	}
	return o1
}

func (o Object) ToPrimitive() Object {
	if o == nil {
		return nil
	}

	toPrimitiveValue := func(v any) any {
		if v == nil {
			return v
		}
		switch v1 := v.(type) {
		case []any, Array, []string, []int:
			return MarshalStringDef(v1, "[]")
		case map[string]any, Object, map[string]string, map[string]int:
			return MarshalStringDef(v1, "{}")
		default:
			return v
		}
	}

	o1 := Object{}
	for k, v := range o {
		o1[k] = toPrimitiveValue(v)
	}
	return o1
}
