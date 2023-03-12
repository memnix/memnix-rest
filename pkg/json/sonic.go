package json

import (
	"github.com/bytedance/sonic"
)

type SonicJSON struct{}

func (*SonicJSON) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func (*SonicJSON) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
