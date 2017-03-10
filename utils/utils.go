package utils

import (
	"golang.org/x/crypto/bcrypt"
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

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}
