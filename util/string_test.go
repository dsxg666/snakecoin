package util

import "testing"

func TestStringIs0ToN(t *testing.T) {
	tests := []struct {
		str string
		res bool
	}{
		{"0", true},
		{"9", true},
		{"10", false},
		{"-1", false},
		{"asd", false},
	}
	for _, tt := range tests {
		if actual := StringIs0ToN(tt.str, 10); actual != tt.res {
			t.Errorf("StringIs0ToN(%s, 10) expected %t, but got %t", tt.str, tt.res, actual)
		}
	}
}

func TestStringIs1ToN(t *testing.T) {
	tests := []struct {
		str string
		res bool
	}{
		{"0", false},
		{"1", true},
		{"10", true},
		{"-1", false},
		{"asd", false},
	}
	for _, tt := range tests {
		if actual := StringIs1ToN(tt.str, 10); actual != tt.res {
			t.Errorf("StringIs0ToN(%s, 10) expected %t, but got %t", tt.str, tt.res, actual)
		}
	}
}
