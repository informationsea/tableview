package main

import (
	"testing";
	"encoding/csv"
	"os"
)

func CheckTestData1(table Table, t *testing.T) {
	if table.GetRow(0)[0] != "Header 1" {
		t.Error("Invalid data at 0,0")
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

	if table.GetMaxColumn() != 3 {
		t.Errorf("Invalid max column: %d", table.GetMaxColumn())
	}

	if table.GetMaxLine() != 72 {
		t.Errorf("Invalid max line: %d", table.GetMaxLine())
	}
}

func TestCreateTable(t *testing.T) {
	input, err1 := os.Open("test1.csv")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer input.Close()

	reader := csv.NewReader(input)
	
	table, err2 := CreateTable(reader)

	if err2 != nil {
		t.Errorf("cannot load data: %s", err2)
		return
	}

	CheckTestData1(table, t)
}

func TestCreateTable2(t *testing.T) {
	input, err1 := os.Open("test1.txt")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer input.Close()

	reader := csv.NewReader(input)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	
	table, err2 := CreateTable(reader)

	if err2 != nil {
		t.Errorf("cannot load data: %s", err2)
		return
	}

	CheckTestData1(table, t)
}

