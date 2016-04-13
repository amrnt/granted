package granted

import (
	"reflect"
	"strings"

	"github.com/mohae/utilitybelt/deepcopy"
)

// FilterToInterface ...
func (a *Authorize) FilterToInterface(i interface{}) interface{} {
	o := deepcopy.Iface(i)
	a.filter(o)
	return o
}

// filter ...
func (a *Authorize) filter(o interface{}) {
	v := reflect.ValueOf(o)

	if v.IsNil() {
		return
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		panic("should be a reflect.Ptr")
	}

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
		if f.CanSet() && !a.hasAccess(v, i) {
			f.Set(reflect.Zero(f.Type()))
		}
	}
}
