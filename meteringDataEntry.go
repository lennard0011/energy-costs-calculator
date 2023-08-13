package main

import "time"

type MeteringDataEntry struct {
	Time    time.Time
	Value   float64 //  The measured energy reading in Wh
	Missing bool    // when constructing and not providing the Missing, it defaults to false.
}

func (myMeterinDataEntry MeteringDataEntry) Weekday() time.Weekday {
	return myMeterinDataEntry.Time.Weekday()
}

func (myMeterinDataEntry MeteringDataEntry) MinutesSinceMidnight() int {
	return myMeterinDataEntry.Time.Hour()*60 + myMeterinDataEntry.Time.Minute()
}
