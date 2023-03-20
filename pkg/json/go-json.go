package json

import (
	"github.com/goccy/go-json"
)

// GoJSON is a wrapper for goccy/go-json.
type GoJSON struct{}

// Marshal marshals the given value to JSON.
func (*GoJSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals the given JSON to the given value.
func (*GoJSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
