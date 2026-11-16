package main

import (
	"testing"
	"time"
)

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name    string
		value   *string
		want    time.Duration
		wantErr bool
	}{
		{"unset uses default", nil, 30 * time.Minute, false},
		{"unit", ptr("45m"), 45 * time.Minute, false},
		{"seconds suffix", ptr("1209600s"), 1209600 * time.Second, false},
		{"no unit is seconds", ptr("1800"), 1800 * time.Second, false},
		{"garbage", ptr("soon"), 0, true},
		{"zero", ptr("0"), 0, true},
		{"negative", ptr("-5m"), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != nil {
				t.Setenv("TEST_DURATION", *tt.value)
			}
			got, err := getDuration("TEST_DURATION", 30*time.Minute)
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("got %v, %v; want %v, err=%v", got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	if defaultTTL != 14*24*time.Hour || defaultCheckInterval != 30*time.Minute {
		t.Fatal("defaults differ from README")
	}
}

func ptr(s string) *string { return &s }
