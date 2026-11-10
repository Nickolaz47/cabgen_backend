package models

import "testing"

func TestParsePlatformStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected PlatformStatus
	}{
		{"ok", PlatformStatusOk},
		{"platformInTesting", PlatformStatusPlatformInTesting},
		{"maintenance", PlatformStatusMaintenance},
		{"adminOnly", PlatformStatusAdminOnly},
		{"analysisDisabled", PlatformStatusAnalysisDisabled},
		{"", PlatformStatusOk},
		{"invalid", PlatformStatusOk},
		{"OK", PlatformStatusOk},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParsePlatformStatus(tt.input); got != tt.expected {
				t.Fatalf("ParsePlatformStatus(%q) = %q, want %q",
					tt.input, got, tt.expected)
			}
		})
	}
}
