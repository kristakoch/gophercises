package main

import (
	"fmt"
	"reflect"
	"testing"
)

// todo: try using structs
// this will remove lines of code, make it more
// readable, and allow specific fail messages
func TestNormalizer(t *testing.T) {
	testCases := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892}",
	}

	want := []string{
		"1234567890",
		"1234567891",
		"1234567892",
		"1234567893",
		"1234567894",
		"1234567890",
		"1234567892",
		"1234567892",
	}

	var got []string
	for _, pn := range testCases {
		res, err := Normalize(pn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		got = append(got, res)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}

}
