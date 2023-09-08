package utils

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	base10  = 10
	bitSize = 32
)

// ConvertStrToUInt converts a string to an unsigned integer
func ConvertStrToUInt(str string) (uint, error) {
	number, err := strconv.ParseUint(str, base10, bitSize)
	if err != nil {
		otelzap.L().Error("Error while converting string to uint", zap.Error(err))
		return 0, errors.New("Error while converting string to uint")
	}
	return uint(number), nil
}

// ConvertUIntToStr converts an unsigned integer to a string
func ConvertUIntToStr(number uint) string {
	return strconv.FormatUint(uint64(number), base10)
}

// ConvertStrToInt converts a string to an integer
func ConvertStrToInt(str string) (int, error) {
	number, err := strconv.ParseInt(str, base10, bitSize)
	if err != nil {
		otelzap.L().Error("Error while converting string to int", zap.Error(err))
		return 0, errors.New("Error while converting string to int")
	}
	return int(number), nil
}
