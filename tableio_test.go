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
	"encoding/csv"
	"os"
	"testing"
)

func CheckTestData1(table Table, t *testing.T) {
	row, err := table.GetRow(0)
	if row[0] != "Header 1" {
		t.Error("Invalid data at 0,0: %s", row[0])
	}

	row, err = table.GetRow(0)
    if row[1] != "Header 2" {
		t.Error("Invalid data at 0,1")
	}

	row, err = table.GetRow(0)
    if row[2] != "Header 3" {
		t.Error("Invalid data at 0,2")
	}

	row, err = table.GetRow(1)
    if row[0] != "1" {
		t.Error("Invalid data at 1,0")
	}

	row, err = table.GetRow(1)
    if row[1] != "2" {
		t.Error("Invalid data at 1,1")
	}

	row, err = table.GetRow(1)
    if row[2] != "3" {
		t.Error("Invalid data at 1,2")
	}
	row, err = table.GetRow(2)
    if row[0] != "12" {
		t.Error("Invalid data at 2,0")
	}

	row, err = table.GetRow(2)
    if row[1] != "23" {
		t.Error("Invalid data at 2,1")
	}

	row, err = table.GetRow(2)
    if row[2] != "34" {
		t.Error("Invalid data at 2,2")
	}

	row, err = table.GetRow(20)
    if row[0] != "1" {
		t.Error("Invalid data at 20,0")
	}

	row, err = table.GetRow(20)
    if row[1] != "2" {
		t.Error("Invalid data at 20,1")
	}

	row, err = table.GetRow(20)
    if row[2] != "3" {
		t.Error("Invalid data at 20,2")
	}

	row, err = table.GetRow(71)
    if row[0] != "E" {
		t.Error("Invalid data at 72,0")
	}

    if row[1] != "N" {
		t.Error("Invalid data at 72,1")
	}

    if row[2] != "D" {
		t.Error("Invalid data at 72,2")
	}
 
	table.LoadAll(1000)
	count, err := table.GetLineCountIfAvailable()

	if err != nil {
		t.Error("Cannot get line count")
	}

	if count != 72 {
		t.Errorf("Invalid max line: %d", count)
	}
}

func CheckTestData4(table Table, t *testing.T) {
	row, err := table.GetRow(0)
    if row[0] != "Header 1" {
		t.Errorf("Invalid data at 0,0: %s", row[2])
	}

	row, err = table.GetRow(0)
    if row[1] != "Header 2" {
		t.Errorf("Invalid data at 0,1: %s", row[2])
	}

	row, err = table.GetRow(0)
    if row[2] != "Header 3" {
		t.Errorf("Invalid data at 0,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("invalid number of columns in first line: %d", row)
	}

	row, err = table.GetRow(1)
    if row[0] != "1" {
		t.Errorf("Invalid data at 1,0: %s", row[0])
	}

	row, err = table.GetRow(1)
    if row[1] != "2" {
		t.Errorf("Invalid data at 1,1: %s", row[1])
	}

	row, err = table.GetRow(1)
    if row[2] != "3" {
		t.Errorf("Invalid data at 1,2: %s", row[2])
	}

	if len(row) != 4 {
		t.Error("invalid number of columns in second line")
	}

	row, err = table.GetRow(2)
    if row[0] != "12" {
		t.Error("Invalid data at 2,0")
	}

	row, err = table.GetRow(2)
    if row[1] != "23" {
		t.Error("Invalid data at 2,1")
	}

	row, err = table.GetRow(2)
    if row[2] != "34" {
		t.Error("Invalid data at 2,2")
	}

	if len(row) != 3 {
		t.Error("invalid number of columns in third line")
	}

	row, err = table.GetRow(20)
    if row[0] != "1" {
		t.Error("Invalid data at 20,0")
	}

	row, err = table.GetRow(20)
    if row[1] != "2" {
		t.Error("Invalid data at 20,1")
	}

	row, err = table.GetRow(20)
    if row[2] != "3" {
		t.Error("Invalid data at 20,2")
	}

	row, err = table.GetRow(71)
    if row[0] != "E" {
		t.Error("Invalid data at 72,0")
	}

	row, err = table.GetRow(71)
    if row[1] != "N" {
		t.Error("Invalid data at 72,1")
	}

	row, err = table.GetRow(71)
    if row[2] != "D" {
		t.Error("Invalid data at 72,2")
	}

	table.LoadAll(1000)
	count, err := table.GetLineCountIfAvailable()

	if err != nil {
		t.Error("Cannot get line count")
	}

	if count != 72 {
		t.Errorf("Invalid max line: %d", count)
	}
}

func TestCreateTable(t *testing.T) {
	data, err1 := LoadTableFromFile("testdata/test1.csv", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable2(t *testing.T) {
	data, err1 := LoadTableFromFile("testdata/test1.txt", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable3(t *testing.T) {
	data, err1 := LoadTableFromFile("testdata/test1.xlsx", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable4(t *testing.T) {
	data, err1 := LoadTableFromFile("testdata/test4.tsv", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData4(data, t)
}

func TestCreateTable5(t *testing.T) {
	data, err1 := LoadTableFromFile("testdata/test3.txt", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData3(data, t)
}

func TestPartialTable1(t *testing.T) {
	input, err := os.Open("testdata/test3.txt")

	if err != nil {
		t.Errorf("Cannot open file: %s", err)
		return
	}

	table := CreatePartialTable(input, ParseTSVRecord)

	defer table.Close()
	CheckTestData3(table, t)
}

func TestPartialTable1CSV(t *testing.T) {
	input, err := os.Open("testdata/test1.csv")

	if err != nil {
		t.Errorf("Cannot open file: %s", err)
		return
	}

	csvReader := csv.NewReader(input)
	table := CreatePartialCSV(csvReader)

	defer table.Close()
	CheckTestData1(table, t)
}

func TestPartialTable2(t *testing.T) {
	input, err := os.Open("testdata/test1.txt")

	if err != nil {
		t.Errorf("Cannot open file: %s", err)
		return
	}

	table := CreatePartialTable(input, ParseTSVRecord)

	defer table.Close()
	CheckTestData1(table, t)
}

func CheckTestData3(table Table, t *testing.T) {
	row, err := table.GetRow(0)
    if row[0] != "日本語のヘッダ" {
		t.Errorf("Invalid data at 0,0: %s", row[0])
	}

	row, err = table.GetRow(0)
    if row[1] != "English Header" {
		t.Errorf("Invalid data at 0,1: %s", row[1])
	}

	row, err = table.GetRow(0)
    if row[2] != "hoge" {
		t.Errorf("Invalid data at 0,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", row)
	}

	row, err = table.GetRow(1)
    if row[0] != "1" {
		t.Errorf("Invalid data at 1,0: %s", row[0])
	}

	row, err = table.GetRow(1)
    if row[1] != "2" {
		t.Errorf("Invalid data at 1,1: %s", row[1])
	}

	row, err = table.GetRow(1)
    if row[2] != "3" {
		t.Errorf("Invalid data at 1,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", row)
	}

	row, err = table.GetRow(2)
    if row[0] != "2" {
		t.Errorf("Invalid data at 2,0: %s", row[0])
	}

	row, err = table.GetRow(2)
    if row[1] != "3" {
		t.Errorf("Invalid data at 2,1: %s", row[1])
	}

	row, err = table.GetRow(2)
    if row[2] != "4" {
		t.Errorf("Invalid data at 2,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", row)
	}

	row, err = table.GetRow(3)
    if row[0] != "3" {
		t.Errorf("Invalid data at 3,0: %s", row[0])
	}

	row, err = table.GetRow(3)
    if row[1] != "4" {
		t.Errorf("Invalid data at 3,1: %s", row[1])
	}

	row, err = table.GetRow(3)
    if row[2] != "5" {
		t.Errorf("Invalid data at 3,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(row))
	}

	row, err = table.GetRow(4)
    if row[0] != "4" {
		t.Errorf("Invalid data at 3,0: %s", row[0])
	}

	row, err = table.GetRow(4)
    if row[1] != "5" {
		t.Errorf("Invalid data at 3,1: %s", row[1])
	}

	row, err = table.GetRow(4)
    if row[2] != "6" {
		t.Errorf("Invalid data at 3,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(row))
	}

	row, err = table.GetRow(5)
    if row[0] != "5" {
		t.Errorf("Invalid data at 3,0: %s", row[0])
	}

	row, err = table.GetRow(5)
    if row[1] != "6" {
		t.Errorf("Invalid data at 3,1: %s", row[1])
	}

	row, err = table.GetRow(5)
    if row[2] != "7" {
		t.Errorf("Invalid data at 3,2: %s", row[2])
	}

	if len(row) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(row))
	}

	table.LoadAll(1000)
	count, err := table.GetLineCountIfAvailable()

	if err != nil {
		t.Error("Cannot get line count")
	}

	if count != 6 {
		t.Errorf("Invalid max line: %d", count)
	}
}

func TestPartialTableEmpty(t *testing.T) {
	input, err := os.Open("testdata/empty.txt")
	
	if err != nil {
		t.Errorf("Cannot open file: %s", err)
		return
	}

	table := CreatePartialTable(input, ParseTSVRecord)
	defer table.Close()

	if table.GetLoadedLineCount() != 0 {
		t.Errorf("Too many loaded line", table.GetLoadedLineCount())
	}
	
	_, err = table.GetRow(0)
	if err == nil {
		t.Error("getting empty column should be null")
	}
}
