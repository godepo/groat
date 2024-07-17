package generics

import (
	"reflect"
	"testing"
)

func Injector[T any, D any](t *testing.T, it T, to D, tag string) D {
	t.Helper()
	pType := reflect.ValueOf(it).Type()

	val := reflect.ValueOf(to)
	tof := reflect.TypeOf(to)

	if tof.Kind() != reflect.Ptr {
		val = reflect.ValueOf(&to).Elem()
	}

	for i := range val.Type().NumField() {
		field := val.Field(i)

		if pType == field.Type() && tof.Field(i).Tag.Get("groat") == tag {
			field.Set(reflect.ValueOf(it))
		}
	}
	return to
}
