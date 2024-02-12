package main

import (
	"log"
)

var (
	namesCSV    string = "./names.csv"
	trainsCSV   string = "./trains.csv"
	namesResult string = "./SD PUS/MAIN/RU/PLACE_NAMES_RESULT.DAT"

	puses = map[string]string{
		// CFG_*
		"ПУС 400": "1",
		"ПУС 401": "2",
		"ПУС 402": "3",
		"ПУС 403": "4",
		"ПУС 406": "5",
		"ПУС 407": "6",
		"ПУС 408": "7",
		"ПУС 409": "8",
	}
)

func main() {
	const op = "main()"

	// get CSV file with names
	csvNames, err := getCSVFile(namesCSV)
	if err != nil {
		log.Printf("%s: open csv file error > %s", op, err)
	}

	defer csvNames.Close()

	// get CSV file with trains
	csvTrains, err := getCSVFile(trainsCSV)
	if err != nil {
		log.Printf("%s: open csv file error > %s", op, err)
	}

	defer csvTrains.Close()

	// get all names from CSV file
	names, err := getCsvLines(csvNames)
	if err != nil {
		log.Printf("%s: get csv lines error > %s", op, err)
	}

	// get all trains from CSV file
	trains, err := getCsvLines(csvTrains)
	if err != nil {
		log.Printf("%s: get csv lines error > %s", op, err)
	}

	// write names with indexes to PLACE_NAMES_RESULT.DAT
	err = writeNamesToResult(names, namesResult)
	if err != nil {
		log.Printf("%s: write names to result error > %s", op, err)
	}

	// write trains with indexes to MSH_A_RESULT.DAT
	err = writeTrainsToResult(trains)
	if err != nil {
		log.Printf("%s: write trains to result error > %s", op, err)
	}
}
