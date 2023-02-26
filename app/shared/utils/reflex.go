package utils

import (
	"reflect"
)

func ExistPropertiesInStruct(v interface{}, s string) bool {
	metaValue := reflect.ValueOf(v).Elem()
	field := metaValue.FieldByName(s)
	return field != (reflect.Value{})
}
