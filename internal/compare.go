package internal

import "reflect"

func CompareType(lhs, rhs reflect.Type) bool {

	if lhs == rhs {
		return true
	}
	if lhs == nil || rhs == nil {
		return false
	}
	return lhs.String() == rhs.String()
}
