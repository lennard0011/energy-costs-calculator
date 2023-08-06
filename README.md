# energy-costs-calculator
Code to calculate the costs of energy consumption based on the input in a CSV format where every row is a 15-minute bucket with the current energy usage.

# Assumptions
* The firs row of the CVS file contains the names of the columns. 
* The timezone of the location is Europe/Amsterdam
* The CSV file is ordered on metering point ids and ascending dates.

# General idea
The general idea of the code is as follows:
1. When run, the main module looks for a CSV file in the input-data folder. 
2. This CSV file is then converted to a measuring data structure per measuring id.
3. The measuring data is filtered for wrong readings.
4. The missing usage between two buckets is imputed by equally dividing the usage of two points.
5. With the complete measuring data, the cost is calculated.

# Architecture
The folder that is being searched for a CSV file is hardcoded   

## Classes
Measuring Data entry
An struct which contains a time (time.Time) and a value (int64). 

Measuring Data

Measuring Data entry

# Questions
1. Is every id is a different household?
2. If we have missing/faulty readings at the end of the history, can we impute it by assuming the usage of the last two points?
3. Can we assume that the first reading for an certain id is always correct? Because we can't calculate the usage.
4. When imputing the wrong values with the linearly, should they be rounded?
