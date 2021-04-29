package main

import "testing"

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
