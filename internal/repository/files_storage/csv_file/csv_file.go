package csv_file

import (
	"encoding/csv"
	"os"
)

type CSVFile struct {
	NameFile string
}

func NewCSVFile(nameFile string) *CSVFile {
	return &CSVFile{NameFile: nameFile}
}

func (C *CSVFile) WriteAll(records [][]string) error {
	err := error(nil)
	file := new(os.File)

	file, err = os.OpenFile(C.NameFile, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	} else {
		if err = os.Truncate(C.NameFile, 0); err != nil {
			panic(err)
		}
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		panic(err)
	}
	return err
}

func (C *CSVFile) ReadAll() ([][]string, error) {
	var records [][]string
	err := error(nil)
	file := new(os.File)

	file, err = os.Open(C.NameFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err = reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return records, err
}
