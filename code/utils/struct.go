package utils

import (
	"reflect"
)

// 根据Key，读取struct value
func GetStructValue[T, V any](stu T, k string) V {
	var value V
	switch reflect.TypeOf(stu).Kind() {
	case reflect.Struct:
		value = reflect.ValueOf(stu).FieldByName(k).Interface().(V)
	case reflect.Ptr:
		isset := reflect.ValueOf(stu).Elem().FieldByName(k)
		if isset.IsValid() {
			value = reflect.ValueOf(stu).Elem().FieldByName(k).Interface().(V)
		}
	}
	return value
}
