package common

import (
	"bufio"
	"log"
	"os"
)

func WriteFile(filename, context string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		log.Panic("Failed to OpenFile: ", err)
	}
	writer := bufio.NewWriter(file)
	writer.WriteString(context)
	writer.Flush()
}

func ReadFile(filename string) []byte {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Panic("Failed to ReadFile: ", err)
	}
	return content
}
