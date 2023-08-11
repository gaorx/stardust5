package sdreflect

import (
	"reflect"
	"slices"
)

func InTypes(t reflect.Type) []reflect.Type {
	numIn := t.NumIn()
	var inTypes []reflect.Type
	for i := 0; i < numIn; i++ {
		inTypes = append(inTypes, t.In(i))
	}
	return inTypes
}

func OutTypes(t reflect.Type) []reflect.Type {
	numOut := t.NumOut()
	var outTypes []reflect.Type
	for i := 0; i < numOut; i++ {
		outTypes = append(outTypes, t.Out(i))
	}
	return outTypes
}

func SplitOutTypes(outTypes []reflect.Type) ([]reflect.Type, bool) {
	n := len(outTypes)
	if n <= 0 {
		return nil, false
	}
	last := outTypes[n-1]
	if last == TErr {
		return slices.Clone(outTypes[0 : n-1]), true
	} else {
		return slices.Clone(outTypes), false
	}
}

func SplitOutValues(outValues []reflect.Value) ([]reflect.Value, reflect.Value) {
	n := len(outValues)
	if n <= 0 {
		return make([]reflect.Value, 0), reflect.Value{}
	}
	last := outValues[n-1]
	if last.Type() == TErr {
		return slices.Clone(outValues[0 : n-1]), outValues[n-1]
	} else {
		return slices.Clone(outValues), reflect.Value{}
	}
}

func MakeFuncIn(inTypes []reflect.Type, f func(inTyp reflect.Type, i int) reflect.Value) []reflect.Value {
	var inVals []reflect.Value
	for i, inTyp := range inTypes {
		outVal := f(inTyp, i)
		inVals = append(inVals, outVal)
	}
	return inVals
}
