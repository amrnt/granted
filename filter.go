package granted

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Filter ...
func Filter(i interface{}) interface{} {
	return Default.filterToInterface(i)
}

// FilterToJSON ...
func FilterToJSON(i interface{}) string {
	j, err := json.Marshal(Default.filterToInterface(i))
	if err != nil {
		return ""
	}

	return string(j)
}

func (a *Authorize) filterToInterface(i interface{}) interface{} {
	a.filter(i)
	return i
}

func (a *Authorize) filter(o interface{}) {
	v := reflect.ValueOf(o)

	if v.IsNil() {
		return
	}

	if v.Kind() != reflect.Ptr {
		// "should be a reflect.Ptr"
		return
	}

	v = v.Elem()
	switch v.Kind() {
	case reflect.Struct:
		a.filterStruct(v)
	case reflect.Map:
		a.filterMapOfStructs(v)
	case reflect.Slice, reflect.Array:
		a.filterSliceOfStructs(v)
	}

	return
}

func (a *Authorize) filterMapOfStructs(v reflect.Value) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		for _, k := range v.MapKeys() {
			f := v.MapIndex(k)
			if f.CanAddr() && f.Kind() != reflect.Ptr {
				f = f.Addr()
			}
			a.filter(f.Interface())
		}
	default:
		a.filter(v.Interface())
	}
}

func (a *Authorize) filterSliceOfStructs(v reflect.Value) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			if f.CanAddr() && f.Kind() != reflect.Ptr {
				f = f.Addr()
			}
			a.filter(f.Interface())
		}
	default:
		a.filter(v.Interface())
	}
}

func (a *Authorize) filterStruct(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		t := v.Type().Field(i).Tag.Get(a.TagName)

		if len(strings.TrimSpace(t)) == 0 {
			continue
		}

		f := v.Field(i)
		switch f.Kind() {
		case reflect.Slice:
			a.filterSliceOfStructs(f)
		case reflect.Map:
			a.filterMapOfStructs(f)
		case reflect.Ptr:
			a.filter(f.Interface())
		case reflect.Struct:
			if f.CanAddr() {
				f = f.Addr()
			}
			a.filter(f.Interface())
		}

		a.filterField(v, t, f)
	}
}

func (a *Authorize) filterField(v reflect.Value, t string, f reflect.Value) {
	ts := strings.Split(t, ",")

	for _, i := range ts {
		// fmt.Println(f.Type())
		if f.CanSet() && !a.canAccess(v, i) {
			f.Set(reflect.Zero(f.Type()))
		}
	}
}
