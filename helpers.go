package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// getCSVFile return csv file by path
func getCSVFile(path string) (*os.File, error) {
	const op = "getCSVFile"

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("%s: get csv file error > %s", op, err)
		return nil, err
	}

	return f, nil
}

// getCsvLines return all records from csv file
func getCsvLines(f *os.File) ([][]string, error) {
	const op = "getCsvLines()"

	// create new reader
	reader := csv.NewReader(f)
	reader.Comma = rune(';')
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("%s: get csv props error > %s", op, err)
		return nil, err
	}

	return records, nil
}

// writeNamesToResult write slice op props to result
func writeNamesToResult(names [][]string, result string) error {
	const op = "writeNamesToResult()"

	f, err := os.OpenFile(result, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("%s: open result file error > %s", op, err)
		return err
	}
	defer f.Close()

	head := "CDA7359AB6408E7F0088CAB68470D5FE" + "\n" + "\n"
	_, err = f.WriteString(head)
	if err != nil {
		log.Printf("%s: write head to file error > %s", op, err)
		return err
	}
	var index string

	l := len(names)
	for i := 0; i < l; i++ {
		switch {
		case i < 9:
			index = "[00" + strconv.Itoa(i+1) + "]" + "\n"
		case i < 99 && i >= 9:
			index = "[0" + strconv.Itoa(i+1) + "]" + "\n"
		case i < 999 && i >= 99:
			index = "[" + strconv.Itoa(i+1) + "]" + "\n"
		case i > 99:
			index = "[" + strconv.Itoa(i+1) + "]" + "\n"
		}

		name := names[i][0] + ",,34567;" + "\n"

		str := index + name

		_, err := f.WriteString(str)
		if err != nil {
			log.Printf("%s: write name to file error > %s", op, err)
			return err
		}
	}

	return nil
}

// write trains
func writeTrainsToResult(trains [][]string) error {
	const op = "writeTrainsToResult()"

	// train
	// 0 = room name
	// 1 = pus name ("ПУС 400")
	// 2 = train number
	// 3 = room index in train

	// CFG_%s == pus index
	// M%s == train number

	var resultFile *os.File

	err := addHeaders()
	if err != nil {
		log.Printf("%s: add headers error > %s", op, err)
		return err
	}

	for _, t := range trains {
		title := t[0]
		pusIndex := puses[t[1]]
		trainNumber := t[2]
		index, _ := strconv.Atoi(t[3])
		resultFilePath := fmt.Sprintf("./SD PUS/CFG_%s/M%s/MSH_A_RESULT.DAT",
			pusIndex, trainNumber)

		if title != "" {
			resultFile, err := os.OpenFile(resultFilePath,
				os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				log.Printf("%s: open result file error > %s", op, err)
				return err
			}

			nameIndex, err := findIndexByTitle(namesResult, title)
			if err != nil {
				log.Printf("%s: find index error > %s", op, err)
			}

			var strIndex string

			switch {
			case index < 10:
				strIndex = "[00" + strconv.Itoa(index) + "]"
			case index < 100 && index >= 10:
				strIndex = "[0" + strconv.Itoa(index) + "]"

			}

			elem := fmt.Sprintf(`
%s
1,,,,,%s,,56789;`, strIndex, nameIndex)

			_, err = resultFile.WriteString(elem)
			if err != nil {
				log.Printf("%s: write train elem to file error > %s", op, err)
				return err
			}
		}
	}

	err = addTails()
	if err != nil {
		log.Printf("%s: write tails to file error > %s", op, err)
		return err
	}

	defer resultFile.Close()

	return nil
}


// addHeaders to every MSH_A_RESULT.DAT file
func addHeaders() error {
	const op = "addHeaders()"

	var i int
	var j int

	// i == CFG_i
	for i = 1; i < 10; i++ {
		// j == Mj
		for j = 1; j < 7; j++ {

			cfg := "CFG_" + strconv.Itoa(i)
			m := "M" + strconv.Itoa(j)
			resultFilePath := fmt.Sprintf("./SD PUS/%s/%s/MSH_A_RESULT.DAT", cfg, m)

			f, err := os.OpenFile(resultFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Printf("%s: open file error > %s", op, err)
				return err
			}

			head := `CDA7359AB6408E7F0088CAB68470D5FE

% серийный номер, код продукта,
% связь А, связь Б, №пом., №зоны
% контрольная сумма (длина 33)`

			f.WriteString(head)
			f.Close()
		}
	}

	return nil
}


// findIndexByTitle return index [index] from PLACE_NAMES_RESULT.DAT by name
func findIndexByTitle(filename string, title string) (string, error) {
	const op = "findIndexByTitle()"

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("%s: open file error > %s", op, err)
		return "", err
	}

	// split by newline
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.Contains(line, title) {
			// check that line above exists
			if i > 0 {
				// return index
				return lines[i-1][1:4], nil
			}

			// if string above is not exists
			log.Printf("%s: name not found in first string > %s", op, err)
			return "", err
		}
	}

	log.Printf("%s: name '%s' is not found in file > %s", op, title, err)
	return "", fmt.Errorf("name '%s' not found in file", title)
}


// addTails add 0,@@@@@@@@,@@@@@@@@,@@@,@@@,@@@,@@@,@@@@@; elements
func addTails() error {
	const op = "addTails()"

	var i int
	var j int

	// i == CFG_i
	for i = 1; i < 10; i++ {
		// j == Mj
		for j = 1; j < 7; j++ {
			cfg := "CFG_" + strconv.Itoa(i)
			m := "M" + strconv.Itoa(j)
			fPath := fmt.Sprintf("./SD PUS/%s/%s/MSH_A_RESULT.DAT", cfg, m)

			content, err := os.ReadFile(fPath)
			if err != nil {
				log.Printf("%s: open file error > %s", op, err)
				return err
			}

			// split by newline
			lines := strings.Split(string(content), "[")
			count := 100
			for _, line := range lines {
				index := "[" + strconv.Itoa(100) + "]"
				if strings.Contains(line, index) {
					continue
				}

				count--
			}

			f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				log.Printf("%s: open file error > %s", op, err)
				return err
			}

			defer f.Close()

			for i := 100 - count; i <= 100; i++ {
				if i < 100 {
					tail := fmt.Sprintf(`
[0%s]
0,@@@@@@@@,@@@@@@@@,@@@,@@@,@@@,@@@,@@@@@;`, strconv.Itoa(i))

					f.WriteString(tail)
				} else {
					tail := `
[100]
0,@@@@@@@@,@@@@@@@@,@@@,@@@,@@@,@@@,@@@@@;`
					f.WriteString(tail)
				}

			}

			f.WriteString("\n% конец файла")
		}
	}

	return nil
}
