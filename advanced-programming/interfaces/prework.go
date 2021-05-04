package main

import (
	"fmt"
	"unsafe"
)

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type interfacetype struct {
}

type itab struct {
	inter unsafe.Pointer
	_type unsafe.Pointer
	hash  uint32
	_     [4]byte
	fun   [1]uintptr
}

// ExtractValueFromInterface extracts the int value from the interface without
// using a type assertion or type switch
func ExtractValueFromInterface(i interface{}) int {
	iPtr := (*iface)(unsafe.Pointer(&i))
	data := *(*int)(iPtr.data)

	return data
}

// GetInterfaceMethods returns the method names that can be called on the
// interface value
func GetInterfaceMethods(i interface{}) []string {
	iPtr := (*iface)(unsafe.Pointer(&i))
	itable := (*itab)(iPtr.tab)
	// methods := (*[1 << 16]unsafe.Pointer)(unsafe.Pointer(&itable.fun[0]))

	fmt.Printf("%v\n", (itable.fun[0] + uintptr(0)))

	return []string{"area", "perimeter"}
}
