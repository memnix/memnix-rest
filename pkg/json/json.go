package json

import "encoding/json"

type Helper interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type JSON struct {
	json Helper
}

func NewJSON(json Helper) *JSON {
	return &JSON{json: json}
}

func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return j.json.Marshal(v)
}

func (j *JSON) Unmarshal(data []byte, v interface{}) error {
	return j.json.Unmarshal(data, v)
}

type NativeJSON struct{}

func (*NativeJSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (*NativeJSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
