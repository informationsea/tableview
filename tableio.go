package main

import (
	"encoding/csv";
	"errors";
	"strings";
	"os"
)

type Table interface {
	GetMaxColumn() int
	GetMaxLine() int
	GetRow(line int) []string
	Close()
}

type SimpleTable struct {
	Data [][]string
	MaxColumn int
	MaxLine int
}

func LoadTableFromFile(filename string, format string) (Table, error) {
	if (format == "auto") {
		if strings.HasSuffix(filename, ".csv") {
			format = "csv"
		} else if strings.HasSuffix(filename, ".txt") {
			format = "tsv"
		} else if strings.HasSuffix(filename, ".tsv") {
			format = "tsv"
		} else if strings.HasSuffix(filename, ".tdf") {
			format = "tsv"
		} else {
			return nil, errors.New("Cannot suggest format")
		}
	} else if (format == "tdf") {
		format = "tsv"
	} else if (format == "csv" || format == "tsv") {
		// ignore
	} else {
		return nil, errors.New("Invalid format")
	}

	inputFile, err3 := os.Open(filename)
	if err3 != nil {
		panic(err3)
	}
	
	defer inputFile.Close()

	if format == "tsv" || format == "csv" {
		csvReader := csv.NewReader(inputFile)
		csvReader.LazyQuotes = true

		if format == "tsv" {
			csvReader.Comma = '\t'
		}
		
		return CreateTable(csvReader)
	}

	return nil, errors.New("Unknown error")
}

func CreateTable(reader *csv.Reader) (Table, error) {
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	
	table := &SimpleTable{data, 0, 0}
	table.MaxLine = len(table.Data)
	table.MaxColumn = 0
	
	for _, v := range table.Data {
		if table.MaxColumn < len(v) {
			table.MaxColumn = len(v)
		}
	}
	
	return table, nil
}

func (v *SimpleTable) GetMaxColumn() int {
	return v.MaxColumn
}

func (v *SimpleTable) GetMaxLine() int {
	return v.MaxLine
}

func (v *SimpleTable) GetRow(line int) []string {
	return v.Data[line]
}

func (v *SimpleTable) Close() {
	// do nothing
}
