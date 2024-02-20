package random

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type Generator struct {
	src rand.Source
}

var (
	instance *Generator //nolint:gochecknoglobals //Singleton
	once     sync.Once  //nolint:gochecknoglobals //Singleton
)

func GetRandomGeneratorInstance() *Generator {
	once.Do(func() {
		instance = &Generator{
			src: rand.NewSource(time.Now().UnixNano()),
		}
	})
	return instance
}

func (r *Generator) GenerateSecretCode(n int) (string, error) {
	sb := strings.Builder{}
	sb.Grow(n)
	// A r.src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, r.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String(), nil
}
