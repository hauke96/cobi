package util

import (
	"fmt"
	"github.com/hauke96/sigolo"
	"reflect"
	"testing"
)

func AssertNil(t *testing.T, value interface{}) {
	if !reflect.DeepEqual(nil, value) {
		sigolo.Errorb(1, "Expect to be 'nil' but was: %+v", value)
		t.Fail()
	}
}

func AssertNotNil(t *testing.T, value interface{}) {
	if nil == value {
		sigolo.Errorb(1, "Expect NOT to be 'nil' but was: %+v", value)
		t.Fail()
	}
}

func AssertError(t *testing.T, expectedMessage string, err error) {
	if expectedMessage != err.Error() {
		sigolo.Errorb(1, "Expected message: %s\nActual error message: %s", expectedMessage, err.Error())
		t.Fail()
	}
}

func AssertEmptyString(t *testing.T, s string) {
	if "" != s {
		sigolo.Errorb(1, "Expected: empty string\nActual  : %s", s)
		t.Fail()
	}
}

func AssertTrue(t *testing.T, b bool) {
	if !b {
		sigolo.Errorb(1, "Expected true but got false")
		t.Fail()
	}
}

func AssertFalse(t *testing.T, b bool) {
	if b {
		sigolo.Errorb(1, "Expected false but got true")
		t.Fail()
	}
}

func AssertEqual[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		sigolo.Errorb(1, "Expected %v but found %v", expected, actual)
		t.Fail()
	}
}

func AssertArrayEqual[T comparable](t *testing.T, expected [][]T, actual [][]T) {
	if len(expected) != len(actual) {
		sigolo.Errorb(1, "Arrays must have the same size in the first dimension, but %d != %d", len(expected), len(actual))
		assertArrayEqualFail(t, expected, actual)
		return
	}

	for x := 0; x < len(expected); x++ {
		if len(expected[x]) != len(actual[x]) {
			sigolo.Errorb(1, "Arrays must have the same size at index %d, but %d != %d", x, len(expected[x]), len(actual[x]))
			assertArrayEqualFail(t, expected, actual)
			return
		}
		for y := 0; y < len(expected[x]); y++ {
			if expected[x][y] != actual[x][y] {
				sigolo.Errorb(1, "Arrays are unequal at [%d, %d]: %d != %d", x, y, len(expected[x]), len(actual[x]))
				assertArrayEqualFail(t, expected, actual)
				return
			}
		}
	}
}

func assertArrayEqualFail[T comparable](t *testing.T, expected [][]T, actual [][]T) {
	t.Fail()
	fmt.Println("Expected: ")
	PrintArray(expected)
	fmt.Println("Actual: ")
	PrintArray(actual)
}
