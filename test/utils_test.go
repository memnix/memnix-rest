package test

import (
	"github.com/memnix/memnixrest/pkg/utils"
	"testing"
)

func TestSendEmail(t *testing.T) {
	type args struct {
		email   string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SendEmail",
			args: args{
				email:   "contact@memnix.app",
				subject: "Test",
				body:    "Test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := utils.SendEmail(tt.args.email, tt.args.subject, tt.args.body); err != nil {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
