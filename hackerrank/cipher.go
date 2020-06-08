package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// caesarCipher returns an encrypted string.
func caesarCipher(s string, k int32) (encrypted string) {
	str := []rune(s)
	for _, ch := range str {
		if unicode.IsLetter(ch) {
			encrypted += encrypt(ch, k)
		} else {
			encrypted += string(ch)
		}
	}
	return encrypted
}

// encrypt returns a character shifted by k positions.
func encrypt(ch int32, k int32) string {
	var lowerBound, upperBound int32
	// Set upper and lower bounds based on case.
	if ch > 64 && ch < 91 {
		lowerBound, upperBound = 65, 90
	} else {
		lowerBound, upperBound = 97, 122
	}
	newPos := (k % 26) + ch
	if newPos > upperBound {
		newPos = lowerBound + (newPos - upperBound - 1)
	}
	return string(newPos)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("CIPHER_OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	n := int32(nTemp)
	_ = n

	s := readLine(reader)

	kTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	k := int32(kTemp)

	result := caesarCipher(s, k)

	fmt.Fprintf(writer, "%s\n", result)

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
