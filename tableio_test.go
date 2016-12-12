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
	if table.GetRow(0)[0] != "Header 1" {
		t.Error("Invalid data at 0,0: %s", table.GetRow(0)[0])
	}

	if table.GetRow(0)[1] != "Header 2" {
		t.Error("Invalid data at 0,1")
	}

	if table.GetRow(0)[2] != "Header 3" {
		t.Error("Invalid data at 0,2")
	}

	if table.GetRow(1)[0] != "1" {
		t.Error("Invalid data at 1,0")
	}

	if table.GetRow(1)[1] != "2" {
		t.Error("Invalid data at 1,1")
	}

	if table.GetRow(1)[2] != "3" {
		t.Error("Invalid data at 1,2")
	}

	if table.GetRow(2)[0] != "12" {
		t.Error("Invalid data at 2,0")
	}

	if table.GetRow(2)[1] != "23" {
		t.Error("Invalid data at 2,1")
	}

	if table.GetRow(2)[2] != "34" {
		t.Error("Invalid data at 2,2")
	}

	if table.GetRow(20)[0] != "1" {
		t.Error("Invalid data at 20,0")
	}

	if table.GetRow(20)[1] != "2" {
		t.Error("Invalid data at 20,1")
	}

	if table.GetRow(20)[2] != "3" {
		t.Error("Invalid data at 20,2")
	}

	if table.GetRow(71)[0] != "E" {
		t.Error("Invalid data at 72,0")
	}

	if table.GetRow(71)[1] != "N" {
		t.Error("Invalid data at 72,1")
	}

	if table.GetRow(71)[2] != "D" {
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
	if table.GetRow(0)[0] != "Header 1" {
		t.Errorf("Invalid data at 0,0: %s", table.GetRow(0)[2])
	}

	if table.GetRow(0)[1] != "Header 2" {
		t.Errorf("Invalid data at 0,1: %s", table.GetRow(0)[2])
	}

	if table.GetRow(0)[2] != "Header 3" {
		t.Errorf("Invalid data at 0,2: %s", table.GetRow(0)[2])
	}

	if len(table.GetRow(0)) != 3 {
		t.Errorf("invalid number of columns in first line: %d", len(table.GetRow(0)))
	}

	if table.GetRow(1)[0] != "1" {
		t.Errorf("Invalid data at 1,0: %s", table.GetRow(1)[0])
	}

	if table.GetRow(1)[1] != "2" {
		t.Errorf("Invalid data at 1,1: %s", table.GetRow(1)[1])
	}

	if table.GetRow(1)[2] != "3" {
		t.Errorf("Invalid data at 1,2: %s", table.GetRow(1)[2])
	}

	if len(table.GetRow(1)) != 4 {
		t.Error("invalid number of columns in second line")
	}

	if table.GetRow(2)[0] != "12" {
		t.Error("Invalid data at 2,0")
	}

	if table.GetRow(2)[1] != "23" {
		t.Error("Invalid data at 2,1")
	}

	if table.GetRow(2)[2] != "34" {
		t.Error("Invalid data at 2,2")
	}

	if len(table.GetRow(2)) != 3 {
		t.Error("invalid number of columns in third line")
	}

	if table.GetRow(20)[0] != "1" {
		t.Error("Invalid data at 20,0")
	}

	if table.GetRow(20)[1] != "2" {
		t.Error("Invalid data at 20,1")
	}

	if table.GetRow(20)[2] != "3" {
		t.Error("Invalid data at 20,2")
	}

	if table.GetRow(71)[0] != "E" {
		t.Error("Invalid data at 72,0")
	}

	if table.GetRow(71)[1] != "N" {
		t.Error("Invalid data at 72,1")
	}

	if table.GetRow(71)[2] != "D" {
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
	if table.GetRow(0)[0] != "日本語のヘッダ" {
		t.Errorf("Invalid data at 0,0: %s", table.GetRow(0)[0])
	}

	if table.GetRow(0)[1] != "English Header" {
		t.Errorf("Invalid data at 0,1: %s", table.GetRow(0)[1])
	}

	if table.GetRow(0)[2] != "hoge" {
		t.Errorf("Invalid data at 0,2: %s", table.GetRow(0)[2])
	}

	if len(table.GetRow(0)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(0)))
	}

	if table.GetRow(1)[0] != "1" {
		t.Errorf("Invalid data at 1,0: %s", table.GetRow(1)[0])
	}

	if table.GetRow(1)[1] != "2" {
		t.Errorf("Invalid data at 1,1: %s", table.GetRow(1)[1])
	}

	if table.GetRow(1)[2] != "3" {
		t.Errorf("Invalid data at 1,2: %s", table.GetRow(1)[2])
	}

	if len(table.GetRow(1)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(1)))
	}

	if table.GetRow(2)[0] != "2" {
		t.Errorf("Invalid data at 2,0: %s", table.GetRow(2)[0])
	}

	if table.GetRow(2)[1] != "3" {
		t.Errorf("Invalid data at 2,1: %s", table.GetRow(2)[1])
	}

	if table.GetRow(2)[2] != "4" {
		t.Errorf("Invalid data at 2,2: %s", table.GetRow(2)[2])
	}

	if len(table.GetRow(2)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(2)))
	}

	if table.GetRow(3)[0] != "3" {
		t.Errorf("Invalid data at 3,0: %s", table.GetRow(3)[0])
	}

	if table.GetRow(3)[1] != "4" {
		t.Errorf("Invalid data at 3,1: %s", table.GetRow(3)[1])
	}

	if table.GetRow(3)[2] != "5" {
		t.Errorf("Invalid data at 3,2: %s", table.GetRow(3)[2])
	}

	if len(table.GetRow(3)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(3)))
	}

	if table.GetRow(4)[0] != "4" {
		t.Errorf("Invalid data at 3,0: %s", table.GetRow(4)[0])
	}

	if table.GetRow(4)[1] != "5" {
		t.Errorf("Invalid data at 3,1: %s", table.GetRow(4)[1])
	}

	if table.GetRow(4)[2] != "6" {
		t.Errorf("Invalid data at 3,2: %s", table.GetRow(4)[2])
	}

	if len(table.GetRow(4)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(4)))
	}

	if table.GetRow(5)[0] != "5" {
		t.Errorf("Invalid data at 3,0: %s", table.GetRow(5)[0])
	}

	if table.GetRow(5)[1] != "6" {
		t.Errorf("Invalid data at 3,1: %s", table.GetRow(5)[1])
	}

	if table.GetRow(5)[2] != "7" {
		t.Errorf("Invalid data at 3,2: %s", table.GetRow(5)[2])
	}

	if len(table.GetRow(5)) != 3 {
		t.Errorf("Invalid number of columns in a cell", len(table.GetRow(5)))
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
