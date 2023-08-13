package main

type WeekdayPriceRate struct {
	Weekday     string
	MinutesFrom int
	MinutesTo   int
	Price       float64 //Price per used Wh
}