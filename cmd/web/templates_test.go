package main

import (
	"testing"
	"time"

	"snippetbox.demien.net/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		want string
		tm   time.Time
	}{
		{
			name: "UTC",
			want: "03 Mar 2026 at 10:45",
			tm:   time.Date(2026, time.March, 3, 10, 45, 0, 0, time.UTC),
		},
		{
			name: "Empty",
			want: "",
			tm:   time.Time{},
		},
		{
			name: "CET",
			want: "03 Mar 2026 at 09:45",
			tm:   time.Date(2026, time.March, 3, 10, 45, 0, 0, time.FixedZone("CET", 1*60*60)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
