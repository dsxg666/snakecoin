package util

import (
	"bufio"
	"log"
	"os"
)

func FileInput(filename, context string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(context)
	writer.Flush()
}

func FileOutput(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
