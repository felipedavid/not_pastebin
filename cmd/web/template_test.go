package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2022, 8, 24, 12, 0, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "24 August at 12:00" {
		t.Errorf("got %q: want %q", hd, "24 August at 12:00")
	}
}
