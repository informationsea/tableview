/*
    tableview: human friendly table viewer
   
    Copyright (C) 2016  OKAMURA, Yasunobu

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors";
	"strings";
	"github.com/tealeg/xlsx";
	"io";
	"bufio";
	"os";
	"encoding/csv"
)

type Table interface {
	GetLineCountIfAvailable() (int, error)
	GetLoadedLineCount() int
	GetRow(int) []string
	LoadAll()
	Load(int)
	Close()
}

type SimpleTable struct {
	Data [][]string
	MaxColumn int
	MaxLine int
}

var TSV_FORMAT = []string{"txt", "tsv", "tdf", "bed", "sam", "gtf", "gff3", "vcf"}

func LoadTableFromFile(filename string, format string) (Table, error) {
	if (format == "auto") {
		if strings.HasSuffix(filename, ".csv") {
			format = "csv"
		} else if strings.HasSuffix(filename, ".xlsx") {
			format = "xlsx"
		} else {
			for _, v := range TSV_FORMAT {
				if strings.HasSuffix(filename, "."+v) {
					format = "tsv"
				}
			}
			
			if format == "auto" {
				return nil, errors.New("Cannot suggest format. Please set -format option.")
			}
		}
	} else if (format == "tdf") {
		format = "tsv"
	} else if (format == "csv" || format == "tsv" || format == "xlsx") {
		// ignore
	} else {
		return nil, errors.New("Invalid format")
	}

	if (format == "xlsx") {
		xlFile, err := xlsx.OpenFile(filename)
		if err != nil {
			panic(err)
		}
		
		sheet := xlFile.Sheets[0]

		data := make([][]string, len(sheet.Rows))
		for ir, row := range sheet.Rows {
			data[ir] = make([]string, len(row.Cells))
			for ic, cell := range row.Cells {
				content, err2 := cell.String()
				if err2 != nil {
					return nil, err2
				}
				
				data[ir][ic] = content
			}
		}
		return CreateTable(data), nil
	}

	inputFile, err3 := os.Open(filename)
	if err3 != nil {
		return nil, err3
	}
	
	defer inputFile.Close()

	if format == "tsv" {
		data, err := LoadTSV(inputFile)
		if err != nil {
			return nil, err
		}
		
		return CreateTable(data), nil
	}

	if format == "csv" {
		csvReader := csv.NewReader(inputFile)

		data, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}
		
		return CreateTable(data), nil
	}

	return nil, errors.New("Unknown error")
}

func LoadTSV(reader io.Reader) ([][]string, error) {
	data := make([][]string, 0)
	
	tabSplit := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i, v := range data {
			if v == '\t' {
				return i+1, data[:i], nil
			}
		}

		return len(data), data, bufio.ErrFinalToken
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		tabScanner := bufio.NewScanner(strings.NewReader(line))
		tabScanner.Split(tabSplit)
		row := make([]string, 0)
		for tabScanner.Scan() {
			cell := tabScanner.Text()
			row = append(row, cell)
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func CreateTable(data [][]string) Table {
	table := &SimpleTable{data, 0, 0}
	table.MaxLine = len(table.Data)
	return table
}

func (v *SimpleTable) GetLineCountIfAvailable() (int, error) {
	return v.MaxLine, nil
}

func (v *SimpleTable) GetLoadedLineCount() int {
	return v.MaxLine
}

func (v *SimpleTable) GetRow(line int) []string {
	return v.Data[line]
}

func (v *SimpleTable) Close() {
	// do nothing
}

func (v *SimpleTable)  LoadAll() {
	// do nothing
}

func (v *SimpleTable)  Load(line int) {
	// do nothing
}
