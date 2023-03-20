package json

import (
	"github.com/goccy/go-json"
)

type GoJson struct{}

func (*GoJson) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (*GoJson) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
