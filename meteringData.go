package main

import "log"

type MeteringData struct {
	ID   string
	Data []MeteringDataEntry
}

func (myMeterinData *MeteringData) addEntry(entry MeteringDataEntry) {
	myMeterinData.Data = append(myMeterinData.Data, entry)
}

func (myMeterinData *MeteringData) flagFaultyValues() {
	// Flag the entries that do not follow the sanity checks
	if len(myMeterinData.Data) < 2 {
		// Can only check the usage if there are at least 2 elements.
		return
	}
	// It is assumed that the first entry is correct.
	previousValue := myMeterinData.Data[0].Value
	for _, meteringDataEntry := range myMeterinData.Data[1:] {
		currentValue := meteringDataEntry.Value
		currentUsage := currentValue - previousValue
		if !usageIsWithinBounds(currentUsage) {
			meteringDataEntry.Missing = true
		}
		previousValue = currentValue
	}
}

func (myMeterinData *MeteringData) linearlyImputeInterval(start, end int) {
	// Linearly impute the values in the interval with the value before and after the interval.

	// Check if start and end are valid indices within the data slice
	if start < 0 || end >= len(myMeterinData.Data) || start >= end {
		log.Println("Invalid start or end index.")
		return
	}

	delta := (myMeterinData.Data[end+1].Value - myMeterinData.Data[start-1].Value) / float64(end-start+2)

	// Apply the linear imputation to each index within the interval
	for i := start; i <= end; i++ {
		myMeterinData.Data[i].Value = myMeterinData.Data[start-1].Value + float64(i-start+1)*delta
	}
}

func (myMeterinData *MeteringData) linearlyForwardImputeInterval(start, end int) {
	// Linearly impute the values in the interval with the value before and after the interval.

	// Check if start and end are valid indices within the data slice
	if start < 0 || end >= len(myMeterinData.Data) || start >= end {
		log.Println("Invalid start or end index.")
		return
	}

	delta := (myMeterinData.Data[start-1].Value - myMeterinData.Data[start-2].Value)

	// Apply the linear imputation to each index within the interval
	for i := start; i <= end; i++ {
		myMeterinData.Data[i].Value = myMeterinData.Data[start-1].Value + float64(i-start+1)*delta
	}
}

func (myMeterinData *MeteringData) imputeMissingValues() {
	// Impute the entries in the data which have the missingDataFiller
	inMissingInterval := false
	missingIntervalIndexStart := -1
	for index, meteringDataEntry := range myMeterinData.Data {
		// loop over all values. If we find a faulty reading we start an interval.
		// Store the previous correct reading and count until it finds a correct reading. When found, inpute the values in between.
		if meteringDataEntry.Missing {
			inMissingInterval = true
			missingIntervalIndexStart = index
		} else {
			if inMissingInterval {
				// here we need to impute the interval of missing values
				inMissingInterval = false
				myMeterinData.linearlyImputeInterval(missingIntervalIndexStart, index-1)
			}
		}
	}
	if inMissingInterval {
		// if end of readings are missing, then inpute last readings linearly from the two readings before.
		myMeterinData.linearlyForwardImputeInterval(missingIntervalIndexStart, len(myMeterinData.Data)-1)
	}
}

func (meteringData MeteringData) calculateCost(weekdayPriceRates WeekdayPriceRateOverview) float64 {
	// find rate that corresponds to meteringdata, then apply calculation.
	// Cant calculate the cost of the first entry as the usage cannot be calculated.
	var totalCosts float64 = 0

	// Need at least two datapoints to calculate usage and costs
	if len(meteringData.Data) < 2 {
		return totalCosts
	}

	for meteringDataEntryIndex, _ := range meteringData.Data[1:] {
		currentMeteringDataEntry := meteringData.Data[meteringDataEntryIndex+1]
		previousMeteringDataEntry := meteringData.Data[meteringDataEntryIndex]

		price := weekdayPriceRates.findPrice(currentMeteringDataEntry.Weekday().String(), currentMeteringDataEntry.MinutesSinceMidnight())
		usage := currentMeteringDataEntry.Value - previousMeteringDataEntry.Value
		cost := price * usage
		totalCosts += cost
	}
	return totalCosts
}
