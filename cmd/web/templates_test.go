package main

import (
	"github.com/felipedavid/not_pastebin/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  string
	}{
		{
			name:  "UTC",
			input: time.Date(2023, 02, 17, 12, 12, 12, 0, time.UTC),
			want:  "17 Feb 2023 at 12:12",
		},
		{
			name:  "Empty object",
			input: time.Time{},
			want:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.input)
			assert.Equal(t, hd, test.want)
		})
	}
}
