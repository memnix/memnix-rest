package json

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

type Mcq struct {
	Answers string `json:"answers" fake:"sentence"`
	Linked  bool   `json:"linked"`
}

// Card is a sample custom struct for testing.
type Card struct {
	Question string `json:"question" fake:"sentence"`
	Answer   string `json:"answer" fake:"sentence"`
	Mcq      Mcq    `json:"mcq"`
	DeckID   uint   `json:"deck_id" fake:"{number:1,100}"`
	McqID    uint   `json:"mcq_id" fake:"{number:1,100}"`
	CardType string `json:"card_type" fake:"{randomstring:[CardTypeMCQOnly,CardTypeBlankMCQ,CardTypeQAProgressive,CardTypeBlankProgressive,CardTypeBlankOnly,CardTypeQAOnly]}"`
}

// Define different sizes of custom struct slices.
var structSliceSizes = []int{100, 200, 500, 1000}

func BenchmarkMarshal(b *testing.B) {
	jsonHelpers := []struct {
		name       string
		jsonHelper Helper
	}{
		{"SonicJSON", &SonicJSON{}},
		{"NativeJSON", &NativeJSON{}},
	}

	for _, size := range structSliceSizes {
		customData := make([]Card, size)
		for i := 0; i < size; i++ {
			err := gofakeit.Struct(&customData[i])
			if err != nil {
				return
			}
		}

		for _, helper := range jsonHelpers {
			b.Run(helper.name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, _ = helper.jsonHelper.Marshal(customData)
				}
			})
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	jsonHelpers := []struct {
		name       string
		jsonHelper Helper
	}{
		{"SonicJSON", &SonicJSON{}},
		{"NativeJSON", &NativeJSON{}},
	}

	for _, size := range structSliceSizes {
		customData := make([]Card, size)
		jsonData, _ := json.Marshal(customData)

		for _, helper := range jsonHelpers {
			b.Run(helper.name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var unmarshaledData []Card
					_ = helper.jsonHelper.Unmarshal(jsonData, &unmarshaledData)
				}
			})
		}
	}
}
