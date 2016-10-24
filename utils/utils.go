package utils

import (
	"reflect"
)

func StructToMap(u interface{}) map[string]string {
	v := reflect.ValueOf(u).Elem()

	m2 := make(map[string]string)

	for i := 0; i < v.NumField(); i++ {
		m2[v.Type().Field(i).Name] = v.Field(i).Interface().(string)
	}

	return m2

}
