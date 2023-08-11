package sdjson

type Array []any

func (a Array) Len() int {
	return len(a)
}

func (a Array) Has(i int) bool {
	return 0 <= i && i < len(a)
}

func (a Array) At(i int) Value {
	if 0 <= i && i < len(a) {
		return V(a[i])
	} else {
		return V(nil)
	}
}

func (a Array) ShadowCopy() Array {
	if a == nil {
		return nil
	}
	a1 := make(Array, 0, len(a))
	for _, elem := range a {
		a1 = append(a1, elem)
	}
	return a1
}
