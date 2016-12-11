/*
    table view: human friendly table viewer
   
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
	"encoding/csv";
	"fmt"
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
	
	//defer inputFile.Close()

	if format == "tsv" {
		return CreateParialTable(inputFile, ParseTSVRecord), nil
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

func CreateTable(data [][]string) Table {
	table := &SimpleTable{data}
	return table
}

func (v *SimpleTable) GetLineCountIfAvailable() (int, error) {
	return len(v.Data), nil
}

func (v *SimpleTable) GetLoadedLineCount() int {
	return len(v.Data)
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

type ParseFinishError struct {}
func (p ParseFinishError) Error() string {
	return "Finished"
}

type ParseRecordFunc func(data string, atEOF bool) ([]string, string, error)

func ParseTSVRecord(data string, atEOF bool) ([]string, string, error) {

	lineEnd := strings.Index(data, "\r\n")

	if lineEnd < 0 {
		lineEnd = strings.Index(data, "\n")
	}

	if lineEnd < 0 {
		if atEOF {
			return strings.Split(data, "\t"), "", ParseFinishError{}
		} else {
			//fmt.Fprintf(os.Stderr, "Line separator is not found", data)
			return nil, data, nil
		}
	} else {
		return strings.Split(data[:lineEnd], "\t"), data[lineEnd+1:], nil
	}
}

type PartialTable struct {
	reader io.Reader
	nextData chan []string
	errChan chan error
	data [][]string
	finish bool
	err error
}

func CreateParialTable(reader io.Reader, parser ParseRecordFunc) *PartialTable {
	nextData := make(chan []string, 100)
	errChan := make(chan error)

	go func() {
		defer close(nextData)
		defer close(errChan)
		scanner := bufio.NewScanner(reader)

		unprocessedData := ""
		for scanner.Scan() {
			unprocessedData += scanner.Text() + "\n"
			//fmt.Fprintf(os.Stderr, "process data: %s\n", unprocessedData)
			data, notProcessed, err := parser(unprocessedData, false)
			unprocessedData = notProcessed
			if data != nil {
				//fmt.Fprintf(os.Stderr, "Next data1 %s\n", data)
				nextData <- data
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error %s", err.Error())
				errChan <- err
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}

		// TODO' loop until finish here
		for len(unprocessedData) > 0 {
			data, notProcessed, err := parser(unprocessedData, true)
			unprocessedData = notProcessed
			if data != nil {
				//fmt.Fprintf(os.Stderr, "Next data2 %s\n", data)
				nextData <- data
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error %s", err.Error())
				errChan <- err
				return
			}
		}
	}()
	
	return &PartialTable{reader, nextData, errChan, make([][]string, 0), false, nil}
}

func (p *PartialTable) loadAvailable() {
	//fmt.Fprintf(os.Stderr, "loadAvailable\n")
	if p.finish {
		return
	}
	
	available := len(p.nextData)
	for i := 0; i < available; i++ {
		p.data = append(p.data, <- p.nextData)
	}
	//fmt.Fprintf(os.Stderr, "loaded 1 %d\n", len(p.data))
}

func (p *PartialTable) GetLineCountIfAvailable() (int, error) {
	//fmt.Fprintf(os.Stderr, "GetLineCountIfAvailable\n")
	p.loadAvailable()
	if p.finish {
		return len(p.data), nil
	} else {
		return -1, errors.New("Not finished")
	}
}

func (p *PartialTable) GetLoadedLineCount() int {
	//fmt.Fprintf(os.Stderr, "GetLoadedLineCount\n")
	p.loadAvailable()
	return len(p.data)
}

func (p *PartialTable) GetRow(line int) []string {
	//fmt.Fprintf(os.Stderr, "GetRow: %d\n", line)
	p.Load(line)
	return p.data[line]
}

func (p *PartialTable) LoadAll() {
	//fmt.Fprintf(os.Stderr, "LoadAll\n")
	for v := range p.nextData {
		p.data = append(p.data, v)
	}
	// fmt.Fprintf(os.Stderr, "loaded 2 %d\n", len(p.data))
	p.finish = true
}

func (p *PartialTable) Load(line int) {
	//fmt.Fprintf(os.Stderr, "Load: %d\n", line)
	if len(p.data) > line {
		//fmt.Fprintf(os.Stderr, "loaded 3 %d/%d\n", line, len(p.data))
		return
	}
	
	for v := range p.nextData {
		p.data = append(p.data, v)
		if len(p.data) > line {
			//fmt.Fprintf(os.Stderr, "loaded 4 %d/%d\n", line, len(p.data))
			return
		}
	}
	//fmt.Fprintf(os.Stderr, "loaded 5 %d/%d\n", line, len(p.data))

	p.finish = true
}

func (p *PartialTable) Close() {
	
}
