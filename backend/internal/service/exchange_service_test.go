package service

import "testing"

func TestResolveExchangeFailReason(t *testing.T) {
	tests := []struct {
		name    string
		reasons []string
		want    string
	}{
		{
			name:    "prefer reject reason",
			reasons: []string{" rejected by risk control ", "fallback failure"},
			want:    "rejected by risk control",
		},
		{
			name:    "fall back to fail reason",
			reasons: []string{"  ", "insufficient balance"},
			want:    "insufficient balance",
		},
		{
			name:    "empty reasons",
			reasons: []string{"", "  "},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveExchangeFailReason(tt.reasons...); got != tt.want {
				t.Fatalf("resolveExchangeFailReason() = %q, want %q", got, tt.want)
			}
		})
	}
}
