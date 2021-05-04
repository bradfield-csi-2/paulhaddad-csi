package main

import (
	"reflect"
	"testing"
)

type Geometry interface {
	area() float64
	perimeter() float64
}

type rectangle struct {
	width, height float64
}

func (r rectangle) area() float64 {
	return r.width * r.height
}

func (r rectangle) perimeter() float64 {
	return 2*r.width + 2*r.height
}

func TestExtractValueFromInterface(t *testing.T) {
	var tests = []struct {
		x    int
		want int
	}{
		{1, 1},
		{10, 10},
	}

	for _, test := range tests {
		var i interface{} = test.x

		got := ExtractValueFromInterface(i)
		if got != test.want {
			t.Errorf("got: %d; want: %d\n", got, test.want)
		}
	}
}

func TestGetInterfaceMethods(t *testing.T) {
	var tests = []struct {
		iFaceValue interface{}
		want       []string
	}{
		{Geometry(rectangle{10.0, 10.0}), []string{"area", "perimeter"}},
	}

	for _, test := range tests {
		got := GetInterfaceMethods(test.iFaceValue)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got: %v; want: %v\n", got, test.want)
		}
	}
}
