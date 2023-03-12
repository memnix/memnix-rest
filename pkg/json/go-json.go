package json

import (
	"github.com/goccy/go-json"
)

type GoJson struct{}

func (g *GoJson) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (g *GoJson) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
