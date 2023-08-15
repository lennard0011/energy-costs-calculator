package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMeteringDataEntry_Weekday(t *testing.T) {
	// Test cases with different weekdays
	tests := []struct {
		timeVal     time.Time
		expectedDay time.Weekday
	}{
		{time.Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC), time.Tuesday},
		{time.Date(2023, time.August, 16, 12, 0, 0, 0, time.UTC), time.Wednesday},
		{time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC), time.Thursday},
		// Add more test cases for other weekdays
	}

	for _, test := range tests {
		entry := MeteringDataEntry{Time: test.timeVal}
		weekday := entry.Weekday()

		assert.Equal(t, weekday, test.expectedDay, "Weekday() should return the correct weekday of the MeteringDataEntry.")
	}
}

func TestMeteringDataEntry_MinutesSinceMidnight(t *testing.T) {
	// Test cases with different times of day
	tests := []struct {
		timeVal         time.Time
		expectedMinutes int
	}{
		{time.Date(2023, time.August, 15, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2023, time.August, 15, 12, 30, 0, 0, time.UTC), 12*60 + 30},
		{time.Date(2023, time.August, 15, 23, 59, 0, 0, time.UTC), 23*60 + 59},
		// Add more test cases for other times
	}

	for _, test := range tests {
		entry := MeteringDataEntry{Time: test.timeVal}
		minutes := entry.MinutesSinceMidnight()

		assert.Equal(t, minutes, test.expectedMinutes, "MinutesSinceMidnight should calculate the minutes since 00:00 of the MeteringDataEntry.")
	}
}
