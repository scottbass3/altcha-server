package config

import (
	"testing"
	"time"
)

func TestExpireDuration(t *testing.T) {
	tests := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{"600", 600 * time.Second, false},
		{"10m", 10 * time.Minute, false},
		{"1h30m", 90 * time.Minute, false},
		{"abc", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		got, err := Config{Expire: tt.in}.ExpireDuration()
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("ExpireDuration(%q) = %v, %v; want %v, err=%v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}
