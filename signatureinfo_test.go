package panobi

import (
	"testing"
	"time"
)

func Test_CalculateSignature(t *testing.T) {
	ki, _ := ParseKey("1234567890123456789012-1234567890123456789012-123")
	now := time.UnixMilli(1672552800000)

	tests := []struct {
		testName          string
		input             string
		wantSignatureInfo SignatureInfo
		wantErr           string
	}{
		{
			testName: "success",
			input:    "Hello, world!",
			wantSignatureInfo: SignatureInfo{
				S:  "v0=8a43a27d205a9d27801d110d1cd627712f2fbf3f123bb6506565917b563def78",
				TS: "1672552800000",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := CalculateSignature([]byte(tt.input), ki, &now)
			if !got.Equals(tt.wantSignatureInfo) {
				t.Errorf("expected key info to be `%v` but got `%v`", tt.wantSignatureInfo, got)
			}
			if !errorIs(tt.wantErr, err) {
				t.Errorf("expected err to be `%s` but got `%v`", tt.wantErr, err)
			}
		})
	}
}
