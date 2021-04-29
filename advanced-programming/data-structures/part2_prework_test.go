package main

import (
	"math"
	"testing"
)

func TestFloat64ToBits(t *testing.T) {
	f := 3.14
	got := Float64ToBits(f)
	want := math.Float64bits(f)

	if got != want {
		t.Errorf("got: %064b; want: %064b\n", got, want)
	}
}

func TestStringSameAddress(t *testing.T) {
	s1 := "hello"
	s2 := "hello world"

	tests := []struct {
		s1, s2 string
		want   bool
	}{
		{s1, s1, true},
		{s1, s2, false},
	}

	for _, test := range tests {
		got := StringSameAddress(test.s1, test.s2)

		if got != test.want {
			t.Errorf("got: %t; want: %t\n", got, test.want)
		}
	}
}

func TestSumOfSlice(t *testing.T) {
	tests := []struct {
		slice []int
		want  int
	}{
		{[]int{1, 2, 3, 4, 5}, int(15)},
	}

	for _, test := range tests {
		got := SumOfSlice(test.slice)

		if got != test.want {
			t.Errorf("got: %d; want: %d\n", got, test.want)
		}
	}
}

// func TestSumOfMap(t *testing.T) {
// 	tests := []struct {
// 		hash map[int]int
// 		want [2]int
// 	}{
// 		{
// 			map[int]int{1: 2, 2: 3, 4: 5},
// 			[2]int{7, 10},
// 		},
// 		{
// 			map[int]int{},
// 			[2]int{0, 0},
// 		},
// 	}
//
// 	for _, test := range tests {
// 		got := SumOfMap(test.hash)
//
// 		if got != test.want {
// 			t.Errorf("got: %v; want: %v", got, test.want)
// 		}
// 	}
// }
