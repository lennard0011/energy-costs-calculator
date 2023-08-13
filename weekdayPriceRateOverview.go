package main

type WeekdayPriceRateOverview struct {
	Rates []WeekdayPriceRate
}

func (myweekdayPriceRateOverview *WeekdayPriceRateOverview) addPriceRate(newWeekdayPriceRate WeekdayPriceRate) {
	myweekdayPriceRateOverview.Rates = append(myweekdayPriceRateOverview.Rates, newWeekdayPriceRate)
}

func (weekdayPriceRates WeekdayPriceRateOverview) findPrice(weekday string, minutesFromMidnight int) float64 {
	var rate float64 = 0
	for _, weekdayPriceRate := range weekdayPriceRates.Rates {
		if weekday == weekdayPriceRate.Weekday && weekdayPriceRate.MinutesFrom <= minutesFromMidnight && minutesFromMidnight <= weekdayPriceRate.MinutesTo {
			rate = weekdayPriceRate.Price
		}
	}
	return rate
}
