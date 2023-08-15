# energy-costs-calculator
Code to calculate the costs of energy consumption based on the input in a CSV format where every row is a 15-minute bucket with the current energy usage.

# How to run
1. Clone the repository.
2. Adjust the pricing.csv according to the electricity prices. Configured per weekday per interval of the day. See the current pricing.csv as a reference. If you want to use a different name for the file, change the value of pricingRatesFileName in main.go.
3. Upload the readings of the meters in the input-data folder in a csv format. See the current test-readings.csv as a reference. The filename is not important, only that is is a csv file. If there are multiple csv files in the folder, then the first one is used.
4. Run the command <code>go run .</code>
5. The calculated costs will be returned in costs.csv.

# Assumptions
* The firs row of the CVS file contains the names of the columns. 
* The timezone of the location is Europe/Amsterdam. This can be adjusted in the main.go
* The CSV file is ordered on metering point ids and ascending dates. 
* There are no missing entries in the file.
* The first reading is correct.
* If a pricerate is not defined in the pricing.csv, it is assumed that the price is 0.

# General idea
The general idea of the code is as follows:
1. When run, the main module looks for a CSV file in the input-data folder. 
2. This CSV file is then converted to a measuring data structure per measuring id.
3. The measuring data is filtered for wrong readings.
4. The missing usage between two buckets is imputed by equally dividing the usage of two points. If the last readings are incorrect, we forward impute this from the two readings before.
5. With the complete measuring data, the cost is calculated.
6. Configuration parameters like usageLowerBound and timezoneLocation can be adjusted in the top of the main.go file.

# Architecture
The data is represented in different classes. The idea is that the codebase is split in different modulair parts where each has a single functionality. The higher-order functions should not care about the complexities of the lower order functions and structs.

The reading and writing of the CSV files is done on a more higher level. Once the data is loaded, it is converted to different classes that handle the logic on the different levels.

Parameters like the usageLowerBound and usageUpperBound can be configured seperatly in the main.go file.

## Classes

### Measuring Data entry
This contains a time (time.Time) and a value (int64). This represents an single readingentry.

### Measuring Data
This represents collection of measuring reading for a single measuring point.

### WeekdayPriceRate
This is a period within the week defined by the weekday and period in the day with a certain price.

### WeekdayPriceRateOverview
This represents a collection of rates which should cover a whole week together used to calculate the prices of usage.

# Next steps
* Define more testcases to fortify the codebase.
* Improve logging to improve user experience.
* Implement different imputation techniques