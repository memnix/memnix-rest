package utils_test

import (
	"testing"

	"github.com/memnix/memnix-rest/pkg/utils"
)

func TestConvertStrToUInt(t *testing.T) {
	tests := []struct {
		str      string
		expected uint
		hasError bool
	}{
		{str: "123", expected: 123, hasError: false},
		{str: "0", expected: 0, hasError: false},
		{str: "4294967295", expected: 4294967295, hasError: false},
		{str: "invalid", expected: 0, hasError: true},
	}

	for _, test := range tests {
		result, err := utils.ConvertStrToUInt(test.str)
		if test.hasError && err == nil {
			t.Errorf("Expected error but got none")
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected ConvertStrToUInt(%s) to return %d, but got %d", test.str, test.expected, result)
		}
	}
}

func TestConvertUIntToStr(t *testing.T) {
	tests := []struct {
		expected string
		number   uint
	}{
		{number: 123, expected: "123"},
		{number: 0, expected: "0"},
		{number: 4294967295, expected: "4294967295"},
	}

	for _, test := range tests {
		result := utils.ConvertUIntToStr(test.number)
		if result != test.expected {
			t.Errorf("Expected ConvertUIntToStr(%d) to return %s, but got %s", test.number, test.expected, result)
		}
	}
}

func TestConvertStrToInt(t *testing.T) {
	tests := []struct {
		str      string
		expected int
		hasError bool
	}{
		{str: "123", expected: 123, hasError: false},
		{str: "0", expected: 0, hasError: false},
		{str: "4294967295", expected: 4294967295, hasError: true},
		{str: "invalid", expected: 0, hasError: true},
	}

	for _, test := range tests {
		result, err := utils.ConvertStrToInt(test.str)
		if test.hasError && err == nil {
			t.Errorf("Expected error but got none")
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected ConvertStrToUInt(%s) to return %d, but got %d", test.str, test.expected, result)
		}
	}
}
