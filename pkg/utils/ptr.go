package utils

import "reflect"

// IsNil checks if interface or its value is nil.
func IsNil(value any) bool {
	if value == nil {
		return true
	}

	refl := reflect.ValueOf(value)

	switch refl.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return refl.IsNil()
	default:
		return false
	}
}
