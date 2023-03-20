package json

import (
	"github.com/bytedance/sonic"
)

// SonicJSON is a wrapper for bytedance/sonic.
type SonicJSON struct{}

// Marshal marshals the given value to JSON.
func (*SonicJSON) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// Unmarshal unmarshals the given JSON to the given value.
func (*SonicJSON) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
