package utils

import "reflect"

// IsNil checks if an interface or its value is nil.
func IsNil(iface any) bool {
	return iface == nil || reflect.ValueOf(iface).IsNil()
}
