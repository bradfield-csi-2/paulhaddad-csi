package main

import (
	"unsafe"
)

type iface struct {
	tab  unsafe.Pointer // this would normally be an itab value, but we don't need it for the exercise
	data unsafe.Pointer
}

// ExtractValueFromInterface extracts the int value from the interface without
// using a type assertion or type switch
func ExtractValueFromInterface(i interface{}) int {
	iPtr := (*iface)(unsafe.Pointer(&i))
	data := *(*int)(iPtr.data)

	return data
}
