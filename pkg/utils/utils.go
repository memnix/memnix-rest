package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ConvertStrToUInt converts a string to an unsigned integer
func ConvertStrToUInt(str string) (uint, error) {
	number, err := strconv.ParseUint(str, config.Base10, config.BitSize)
	if err != nil {
		log.Debug().Err(err).Msgf("Error while converting string to uint: %s", err)
		return 0, errors.New("Error while converting string to uint")
	}
	return uint(number), nil
}

// ConvertUIntToStr converts an unsigned integer to a string
func ConvertUIntToStr(number uint) string {
	return strconv.FormatUint(uint64(number), config.Base10)
}

// ConvertStrToInt converts a string to an integer
func ConvertStrToInt(str string) (int, error) {
	number, err := strconv.ParseInt(str, config.Base10, config.BitSize)
	if err != nil {
		log.Debug().Err(err).Msgf("Error while converting string to int: %s", err)
		return 0, errors.New("Error while converting string to int")
	}
	return int(number), nil
}

// GetExpirationTime returns the expiration time
func GetExpirationTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Hour * config.ExpirationTimeInHours))
}

// GetSecretKey returns the secret key
func GetSecretKey() string {
	return config.EnvHelper.GetEnv("SECRET_KEY")
}
