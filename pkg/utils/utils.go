package utils

import (
	"strconv"
)

const (
	base10  = 10
	bitSize = 32
)

// ConvertStrToUInt converts a string to an unsigned integer.
func ConvertStrToUInt(str string) (uint, error) {
	number, err := strconv.ParseUint(str, base10, bitSize)
	if err != nil {
		return 0, err
	}
	return uint(number), nil
}

// ConvertUIntToStr converts an unsigned integer to a string.
func ConvertUIntToStr(number uint) string {
	return strconv.FormatUint(uint64(number), base10)
}

// ConvertIntToStr converts an integer to a string.
func ConvertIntToStr(number int) string {
	return strconv.FormatInt(int64(number), base10)
}

// ConvertInt32ToStr converts an int32 to a string.
func ConvertInt32ToStr(number int32) string {
	return strconv.FormatInt(int64(number), base10)
}

// ConvertStrToInt converts a string to an integer.
func ConvertStrToInt(str string) (int, error) {
	number, err := strconv.ParseInt(str, base10, bitSize)
	if err != nil {
		return 0, err
	}
	return int(number), nil
}

// ConvertStrToInt32 converts a string to an int32.
func ConvertStrToInt32(str string) (int32, error) {
	number, err := strconv.ParseInt(str, base10, bitSize)
	if err != nil {
		return 0, err
	}
	return int32(number), nil
}
