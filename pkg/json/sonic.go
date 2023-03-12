package json

import (
	"github.com/bytedance/sonic"
)

type SonicJSON struct{}

func (s *SonicJSON) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func (s *SonicJSON) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
