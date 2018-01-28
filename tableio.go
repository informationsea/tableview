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
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type Table interface {
	GetLineCountIfAvailable() (int, error)
	GetLoadedLineCount() int
	GetRow(int) ([]string, error)
	LoadAll(timeout int) bool
	Load(line int)
	Close()
}

type SimpleTable struct {
	Data [][]string
}

var TSV_FORMAT = []string{"txt", "tsv", "tdf", "bed", "sam", "gtf", "gff3", "vcf"}

func LoadTableFromFile(filename string, format string, sheetNum int) (Table, error) {
	if format == "xlsx" || (strings.HasSuffix(filename, ".xlsx") && format == "auto") {
		xlFile, err := xlsx.OpenFile(filename)
		if err != nil {
			panic(err)
		}

		if len(xlFile.Sheets) < sheetNum-1 {
			panic("Invalid sheet number")
		}

		sheet := xlFile.Sheets[sheetNum-1]

		data := make([][]string, len(sheet.Rows))
		for ir, row := range sheet.Rows {
			data[ir] = make([]string, len(row.Cells))
			for ic, cell := range row.Cells {
				content := cell.String()
				//if err2 == nil {
				data[ir][ic] = content
				//} else {
				//	data[ir][ic] = cell.Value
				//}
			}
		}
		return CreateTable(data), nil
	}

	var reader io.Reader
	var err error
	reader, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filename, ".gz") {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		filename = filename[:len(filename)-3]
	}

	if strings.HasSuffix(filename, ".bz2") {
		reader = bzip2.NewReader(reader)
		filename = filename[:len(filename)-4]
	}

	if format == "auto" {
		if strings.HasSuffix(filename, ".csv") {
			format = "csv"
		} else {
			format = "tsv"
		}
	} else if format == "tdf" {
		format = "tsv"
	} else if format == "csv" || format == "tsv" {
		// ignore
	} else {
		return nil, errors.New("Invalid format")
	}

	//defer inputFile.Close()

	if format == "tsv" {
		return CreatePartialTable(reader, ParseTSVRecord), nil
	} else if format == "csv" {
		csvReader := csv.NewReader(reader)
		return CreatePartialCSV(csvReader), nil
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

func (v *SimpleTable) GetRow(line int) ([]string, error) {
	//if len(v.Data[line]) <= line {
	//		return nil, errors.New("index out of range")
	//	}
	return v.Data[line], nil
}

func (v *SimpleTable) Close() {
	// do nothing
}

func (v *SimpleTable) LoadAll(timeout int) bool {
	// do nothing
	return true
}

func (v *SimpleTable) Load(line int) {
	// do nothing
}

type ParseFinishError struct{}

func (p ParseFinishError) Error() string {
	return "Finished"
}

type ParseRecordFunc func(data string, atEOF bool) ([]string, string, error)

func ParseTSVRecord(data string, atEOF bool) ([]string, string, error) {
	lineEnd := strings.Index(data, "\n")

	if lineEnd < 0 {
		if atEOF {
			return strings.Split(data, "\t"), "", ParseFinishError{}
		} else {
			return nil, data, nil
		}
	} else {
		return strings.Split(data[:lineEnd], "\t"), data[lineEnd+1:], nil
	}
}

type PartialTable struct {
	nextData chan []string
	errChan  chan error
	data     [][]string
	finish   bool
	err      error
}

func CreatePartialCSV(reader *csv.Reader) *PartialTable {
	nextData := make(chan []string, 1000)
	errChan := make(chan error)

	go func() {
		defer close(nextData)
		defer close(errChan)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- err
				return
			}
			nextData <- record
		}
	}()

	return &PartialTable{nextData, errChan, make([][]string, 0), false, nil}
}

func CreatePartialTable(reader io.Reader, parser ParseRecordFunc) *PartialTable {
	nextData := make(chan []string, 1000)
	errChan := make(chan error)

	go func() {
		defer close(nextData)
		defer close(errChan)
		scanner := bufio.NewScanner(reader)

		unprocessedData := ""
		for scanner.Scan() {
			unprocessedData += scanner.Text() + "\n"
			data, notProcessed, err := parser(unprocessedData, false)
			unprocessedData = notProcessed
			if data != nil {
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

		for len(unprocessedData) > 0 {
			data, notProcessed, err := parser(unprocessedData, true)
			unprocessedData = notProcessed
			if data != nil {
				nextData <- data
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error %s", err.Error())
				errChan <- err
				return
			}
		}
	}()

	return &PartialTable{nextData, errChan, make([][]string, 0), false, nil}
}

func (p *PartialTable) GetLineCountIfAvailable() (int, error) {
	if p.finish {
		return len(p.data), nil
	} else {
		return -1, errors.New("Not finished")
	}
}

func (p *PartialTable) GetLoadedLineCount() int {
	return len(p.data) + len(p.nextData)
}

func (p *PartialTable) GetRow(line int) ([]string, error) {
	p.Load(line)
	if len(p.data) <= line || line < 0 {
		return nil, errors.New("index out of range line")
	}
	return p.data[line], nil
}

func (p *PartialTable) LoadAll(timeout int) bool {
	start := time.Now()
	for v := range p.nextData {
		p.data = append(p.data, v)
		d := time.Since(start)
		if d.Nanoseconds() > int64(timeout)*1000 {
			return false
		}
	}
	p.finish = true
	return true
}

func (p *PartialTable) Load(line int) {
	if len(p.data) > line {
		return
	}

	for v := range p.nextData {
		p.data = append(p.data, v)
		if len(p.data) > line {
			return
		}
	}
	p.finish = true
}

func (p *PartialTable) Close() {

}
