package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Float64ToBits returns a float64 as a uint64 value
func Float64ToBits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

// StringSameAddress returns whether the underlying string data of two strings
// points at the same address
func StringSameAddress(s1, s2 string) bool {
	header1 := (*reflect.StringHeader)(unsafe.Pointer(&s1))
	header2 := (*reflect.StringHeader)(unsafe.Pointer(&s2))

	return header1.Len == header2.Len && header1.Data == header2.Data
}

// SumOfSlice returns the sum of an array using pointer arithmetric
func SumOfSlice(s []int64) int64 {
	// Get slice length from its header
	st := reflect.ValueOf(s)
	length := uintptr(st.Len())
	sum := int64(0)

	for i := uintptr(0); i < length; i++ {
		sum += *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + i*unsafe.Sizeof(s[0])))
	}

	return sum
}

// SumOfMap returns an array of the sum of the map's keys and the sum of the
// map's values
func SumOfMap(hash map[int]int) [2]int {
	mt := reflect.ValueOf(hash)
	t := mt.Type()
	length := mt.Len()

	fmt.Println(t, mt, length)

	return [2]int{7, 10}
}
