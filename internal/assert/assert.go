package assert

import (
	"fmt"
	"reflect"
)

func AssertMap(v reflect.Value) {
	switch v.Kind() {
	case reflect.Map:
	default:
		panic(fmt.Sprintf("Assertion failed: %v is not a Map", v))
	}
}

func AssertStruct(v reflect.Value) {
	vk := v.Kind()
	switch vk {
	case reflect.Struct:
	case reflect.Pointer:
		AssertStruct(v.Elem())
	default:
		panic(fmt.Sprintf("Assertion failed: %s is not Struct", vk))
	}
}

func AssertType(v reflect.Value, t reflect.Type) {
	if v.Type() == t {
		return
	}

	panic(fmt.Sprintf("Assertion failed: Value %v is not of Type %v", v.Type(), t))
}
