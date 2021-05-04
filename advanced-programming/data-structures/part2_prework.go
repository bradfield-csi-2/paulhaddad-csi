package main

import (
	"unsafe"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

type sliceStruct struct {
	array unsafe.Pointer
	len   int
	cap   int
}

// Float64ToBits returns a float64 as a uint64 value
func Float64ToBits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

// StringSameAddress returns whether the underlying string data of two strings
// points at the same address
func StringSameAddress(s1, s2 string) bool {
	strStruct1 := (*stringStruct)(unsafe.Pointer(&s1))
	strStruct2 := (*stringStruct)(unsafe.Pointer(&s2))

	return strStruct1.str == strStruct2.str
}

// SumOfSlice returns the sum of an array using pointer arithmetric
func SumOfSlice(s []int) int {
	sliceStr := (*sliceStruct)(unsafe.Pointer(&s))
	base := uintptr(sliceStr.array)

	sum := 0
	for i := 0; i < sliceStr.len; i++ {
		offset := uintptr(i) * unsafe.Sizeof(int(0))
		sum += *(*int)(unsafe.Pointer(base + offset))
	}
	return sum
}

// SumOfMap returns an array of the sum of the map's keys and the sum of the
// map's values
// func SumOfMap(hash map[int]int) [2]int {
// 	mt := reflect.ValueOf(hash)
// 	t := mt.Type()
// 	length := mt.Len()
//
// 	fmt.Println(t, mt, length)
//
// 	return [2]int{7, 10}
// }
