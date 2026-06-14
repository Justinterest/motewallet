package kun

import "testing"

func TestPlatformAccountToKUN(t *testing.T) {
	tests := []struct {
		platform string
		want     string
		wantErr  bool
	}{
		{"FUNDING", AccountHK, false},
		{"TRADING", AccountPL, false},
		{"UNKNOWN", "", true},
	}

	for _, tt := range tests {
		got, err := PlatformAccountToKUN(tt.platform)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("platform=%s expected error", tt.platform)
			}
			continue
		}
		if err != nil {
			t.Fatalf("platform=%s unexpected error: %v", tt.platform, err)
		}
		if got != tt.want {
			t.Fatalf("platform=%s got %s want %s", tt.platform, got, tt.want)
		}
	}
}
