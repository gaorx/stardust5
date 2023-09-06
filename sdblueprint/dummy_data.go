package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"reflect"
)

const (
	dummyDataFieldName = "DummyData"
)

type DummyRecord map[string]any

func (record DummyRecord) ToDB(t Table) DummyRecord {
	if record == nil {
		return nil
	}
	dbRecord := DummyRecord{}
	for colId, v := range record {
		col := t.Column(colId)
		if col == nil {
			panic(sderr.NewWith("not found column when fill dummy data", sderr.Attrs{"c": colId, "t": t.Id()}))
		}
		dbRecord[col.NameForDB()] = v
	}
	return dbRecord
}

func (record DummyRecord) Has(colId string) bool {
	_, ok := record[colId]
	return ok
}

func (record DummyRecord) Get(colId string) any {
	v := record[colId]
	return v
}

func (record DummyRecord) Lookup(colId string) (any, bool) {
	v, ok := record[colId]
	return v, ok
}

func getDummyData(session any, sv reflect.Value, st structType, cols []column) ([]DummyRecord, error) {
	var dummyData any

	// get dummy data in struct field or method
	if dummyDataVal := sv.FieldByName(dummyDataFieldName); dummyDataVal.IsValid() {
		dummyData = dummyDataVal.Interface()
	} else if dummyDataMethod := sv.MethodByName(dummyDataFieldName); dummyDataMethod.IsValid() {
		if dummyDataMethod.Type().NumOut() != 1 {
			return nil, sderr.New("illegal signature for DummyData method")
		}
		switch dummyDataMethod.Type().NumIn() {
		case 0:
			dummyData = dummyDataMethod.Call([]reflect.Value{})[0].Interface()
		case 1:
			dummyData = dummyDataMethod.Call([]reflect.Value{sdreflect.ValueOf(session)})[0].Interface()
		default:
			return nil, sderr.New("illegal signature for DummyData method")
		}
	}
	if dummyData == nil {
		return nil, nil
	}

	// if dummy data is a callback, call it
	if ddv := reflect.ValueOf(dummyData); ddv.Kind() == reflect.Func {
		if ddv.Type().NumOut() != 1 {
			return nil, sderr.New("illegal signature for DummyData callback")
		}
		switch ddv.Type().NumIn() {
		case 0:
			dummyData = ddv.Call([]reflect.Value{})[0].Interface()
		case 1:
			dummyData = ddv.Call([]reflect.Value{sdreflect.ValueOf(session)})[0].Interface()
		default:
			return nil, sderr.New("illegal signature for DummyData callback")
		}
	}
	if dummyData == nil {
		return nil, nil
	}

	// dummy data must be a slice/array
	if ddv := reflect.ValueOf(dummyData); ddv.Kind() != reflect.Slice && ddv.Kind() != reflect.Array {
		return nil, sderr.New("dummy data must be a slice/array")
	}

	// normalize
	ddv := reflect.ValueOf(dummyData)
	var records []DummyRecord
	n := ddv.Len()
	for i := 0; i < n; i++ {
		record, ok := toDummyRecord(ddv.Index(i), cols)
		if !ok {
			return nil, sderr.New("illegal dummy record")
		}
		records = append(records, record)
	}
	return records, nil
}

func toDummyRecord(recordVal reflect.Value, cols []column) (DummyRecord, bool) {
	if recordVal.Kind() == reflect.Pointer {
		return toDummyRecord(recordVal.Elem(), cols)
	} else if recordVal.Kind() == reflect.Struct {
		r := DummyRecord{}
		for _, c := range cols {
			colId := c.id
			colVal := recordVal.FieldByName(colId)
			if colVal.IsValid() {
				r[colId] = colVal.Interface()
			} else {
				r[colId] = c.Default()
			}
		}
		return r, true
	} else if recordVal.Kind() == reflect.Map {
		if recordVal.Type().Key() != sdreflect.TString {
			return nil, false
		}
		r := DummyRecord{}
		iter := recordVal.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			r[k.Interface().(string)] = v.Interface()
		}
		for _, c := range cols {
			_, ok := r[c.id]
			if !ok {
				r[c.id] = c.Zero()
			}
		}
		return r, true
	} else {
		return nil, false
	}
}
