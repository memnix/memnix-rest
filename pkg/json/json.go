package json

import "encoding/json"

// Helper is an interface for JSON helper.
type Helper interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// JSON is a wrapper for JSON helper.
type JSON struct {
	json Helper
}

// NewJSON returns a new JSON.
func NewJSON(json Helper) *JSON {
	return &JSON{json: json}
}

// Marshal marshals the given value to JSON.
func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}

// Unmarshal unmarshals the given JSON to the given value.
func (j *JSON) Unmarshal(data []byte, v interface{}) error {
	return j.json.Unmarshal(data, v)
}

// NativeJSON is a wrapper for encoding/json.
type NativeJSON struct{}

// Marshal marshals the given value to JSON.
func (*NativeJSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals the given JSON to the given value.
func (*NativeJSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
