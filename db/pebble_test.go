package db

import (
	"bytes"
	"testing"
)

func TestDB(t *testing.T) {
	db := GetDB("./test")
	tests := []struct{ k, v, res string }{
		{"hello", "world", "world"},
		{"dsxg", "handsome", "handsome"},
	}
	for _, tt := range tests {
		Set([]byte(tt.k), []byte(tt.v), db)
		if actual := Get([]byte(tt.k), db); bytes.Compare(actual, []byte(tt.res)) != 0 {
			t.Error("unexpected occur")
		}
	}
	CloseDB(db)
}
