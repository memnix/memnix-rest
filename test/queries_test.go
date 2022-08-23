package test

import (
	"github.com/memnix/memnixrest/app/queries"
	"reflect"
	"testing"
)

func TestFetchTodayCard(t *testing.T) {
	tests := []struct {
		name   string
		userID uint
		want   bool
	}{
		{
			name:   "fetch today card",
			userID: 6,
			want:   true,
		},
	}

	_, err := Setup()
	if err != nil {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := queries.FetchTodayCard(tt.userID); !reflect.DeepEqual(got.Success, tt.want) {
				t.Errorf("FetchTodayCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
