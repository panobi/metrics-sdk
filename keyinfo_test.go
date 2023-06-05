package panobi

import (
	"testing"
)

func Test_ParseKey(t *testing.T) {
	tests := []struct {
		testName    string
		input       string
		wantKeyInfo KeyInfo
		wantErr     string
	}{
		{
			testName:    "empty input",
			input:       "",
			wantKeyInfo: KeyInfo{},
			wantErr:     "invalid key",
		},
		{
			testName:    "missing parts",
			input:       "abc",
			wantKeyInfo: KeyInfo{},
			wantErr:     "invalid key",
		},
		{
			testName:    "WID too short",
			input:       "abc-",
			wantKeyInfo: KeyInfo{},
			wantErr:     "invalid key",
		},
		{
			testName:    "key too short",
			input:       "1234567890123456789012-",
			wantKeyInfo: KeyInfo{},
			wantErr:     "invalid key",
		},
		{
			testName: "success",
			input:    "1234567890123456789012-1234567890123456789012-def",
			wantKeyInfo: KeyInfo{
				WorkspaceID: "1234567890123456789012",
				ExternalID:  "1234567890123456789012",
				K:           "def",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := ParseKey(tt.input)
			if !got.Equals(tt.wantKeyInfo) {
				t.Errorf("expected key info to be `%v` but got `%v`", tt.wantKeyInfo, got)
			}
			if !errorIs(tt.wantErr, err) {
				t.Errorf("expected err to be `%s` but got `%v`", tt.wantErr, err)
			}
		})
	}
}
