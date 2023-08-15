package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const inputfoldername string = "input-data"
const pricingRatesFileName string = "pricing.csv"
const missingDataFiller string = "missing value"
const timezoneLocation string = "Europe/Amsterdam"
const usageLowerBound float64 = 0
const usageUpperBound float64 = 100

func usageIsWithinBounds(usage float64) bool {
	// Sanity checks a usage. It checks if the usage between not be less then 0 and more than 100
	return usage >= usageLowerBound && usage <= usageUpperBound
}

func importPricingRates() (WeekdayPriceRateOverview, error) {
	pricingRates := WeekdayPriceRateOverview{}

	// Open the CSV file
	file, err := os.Open(pricingRatesFileName)
	if err != nil {
		return pricingRates, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		return pricingRates, err
	}

	for _, record := range records[1:] {
		weekday := record[0]
		minutesFrom, _ := strconv.Atoi(record[1])
		minutesTo, _ := strconv.Atoi(record[2])
		price, _ := strconv.ParseFloat(record[3], 64)

		// Create a new weekdayPriceRate struct and append it to the slice
		weekdayPriceRate := WeekdayPriceRate{
			Weekday:     weekday,
			MinutesFrom: minutesFrom,
			MinutesTo:   minutesTo,
			Price:       price,
		}
		pricingRates.addPriceRate(weekdayPriceRate)
	}

	return pricingRates, nil
}

func findCSVInFolder() string {
	// Get a list of all files in the folder
	files, err := os.ReadDir(inputfoldername)
	if err != nil {
		panic(err)
	}

	// Iterate through the files and find the first CSV file
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// Check if the file has a .csv extension (case-insensitive)
		if strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			// Found a CSV file
			csvFilePath := filepath.Join(inputfoldername, file.Name())
			return csvFilePath // Exit the loop after finding the first CSV file
		}
	}

	// No CSV file found in the input folder
	return ""
}

func readCSVtoMeteringData(csvFilePath string, timezoneName string) (map[string]MeteringData, error) {
	var meteringDatas = make(map[string]MeteringData)

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return meteringDatas, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	dataRows, err := reader.ReadAll()
	if err != nil {
		return meteringDatas, err
	}
	// First row consists of the headers of the CSV
	records := dataRows[1:]

	loc, err := time.LoadLocation(timezoneName)
	if err != nil {
		return meteringDatas, err
	}

	// Process each record
	var currentId string
	var currentMeteringData MeteringData
	// Start from the second row as the first are the headers of the CSV
	for _, record := range records {
		id := record[0]
		value := record[1]
		timerecord := record[2]

		if currentId != id {
			if currentId != "" {
				meteringDatas[currentId] = currentMeteringData
			}
			currentMeteringData = MeteringData{}
			currentMeteringData.ID = id
			currentId = id
		}

		// Parse the data to the appropriate types
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Println("Error converting reading to int:", err)
			continue
		}
		timeInt, err := strconv.ParseInt(timerecord, 10, 64)
		myTime := time.Unix(timeInt, 0).In(loc)
		if err != nil {
			fmt.Println("Error converting createdAt to uint64:", err)
			continue
		}

		// Create a new MeteringData struct and append it to the slice
		entry := MeteringDataEntry{
			Time:  myTime,
			Value: valueFloat,
		}
		currentMeteringData.addEntry(entry)
	}

	meteringDatas[currentId] = currentMeteringData

	// Now we have the meteringDatas in a map where the id's are a key
	return meteringDatas, nil
}

func main() {
	// In a specific folder take the first CSV you find
	csvFilePath := findCSVInFolder()
	if csvFilePath == "" {
		log.Println("No CSV file for reading metering found in input folder.")
		return
	}

	// Per Metering ID, Transform the CSV to a map of MeteringData
	meterinDataMap, err := readCSVtoMeteringData(csvFilePath, timezoneLocation)
	if err != nil {
		log.Fatalf("Failed reading metering data: %s", err)
	}

	energyCosts := make(map[string]float64)

	energyPriceRates, err := importPricingRates()
	if err != nil {
		log.Fatalf("Failed reading price data: %s", err)
	}

	for id, meteringData := range meterinDataMap {
		// Filter out faulty readings
		meteringData.flagFaultyValues()
		// Impute the missing readings
		meteringData.imputeMissingValues()
		// Calculate the costs of the energy used
		energyCosts[id] = meteringData.calculateCost(energyPriceRates)
	}

	csvFile, err := os.Create("costs.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	headers := []string{"id", "cost"}
	_ = csvwriter.Write(headers)
	for id, cost := range energyCosts {
		_ = csvwriter.Write([]string{id, strconv.FormatFloat(cost, 'f', -1, 64)})
	}
	csvwriter.Flush()

	csvFile.Close()
}
