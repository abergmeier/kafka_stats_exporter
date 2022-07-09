package assert

import (
	"fmt"
	"reflect"
)

func AssertType(v reflect.Value, t reflect.Type) {
	if v.Type() == t {
		return
	}

	panic(fmt.Sprintf("Assertion failed: Value %v is not of Type %v", v.Type(), t))
}
