package util

import (
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	FileInput("./test/text.txt", "hello")
	context := FileOutput("./test/text.txt")
	if strings.Compare("hello", context) != 0 {
		t.Error("unexpected occur")
	}
}
