package common

import "testing"

func TestStringIsAllNumber(t *testing.T) {
	tests := []struct {
		num string
		res bool
	}{
		{"0", true},
		{"0.9", true},
		{"100", true},
		{"a", false},
	}
	for _, tt := range tests {
		if actual := StringIsAllNumber(tt.num); actual != tt.res {
			t.Errorf("StringIsAllNumber(%s) expected %t, but got %t", tt.num, tt.res, actual)
		}
	}
}
