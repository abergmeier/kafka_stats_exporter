package gen

import (
	"fmt"
	"reflect"
)

func assertMap(v reflect.Value) {
	switch v.Kind() {
	case reflect.Map:
	default:
		panic(fmt.Sprintf("Assertion failed: %v is not a Map", v))
	}
}

func assertStruct(v reflect.Value) {
	vk := v.Kind()
	switch vk {
	case reflect.Struct:
	case reflect.Pointer:
		assertStruct(v.Elem())
	default:
		panic(fmt.Sprintf("Assertion failed: %s is not Struct", vk))
	}
}
