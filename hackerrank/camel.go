package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Complete the camelcase function below.
func camelcase(s string) int {
	numWords := 1
	if len(s) == 1 {
		return 1
	}
	s = s[1:]
	for _, ch := range s {
		if ch > 64 && ch < 91 {
			numWords++
		}
	}
	return numWords
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("CAMEL_OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	s := readLine(reader)

	result := camelcase(s)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
