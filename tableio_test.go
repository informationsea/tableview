package main

import (
	"testing"
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

func CheckTestData4(table Table, t *testing.T) {
	if table.GetRow(0)[0] != "Header 1" {
		t.Error("Invalid data at 0,0")
	}

	if table.GetRow(0)[1] != "Header 2" {
		t.Error("Invalid data at 0,1")
	}

	if table.GetRow(0)[2] != "Header 3" {
		t.Error("Invalid data at 0,2")
	}

	if len(table.GetRow(0)) != 3 {
		t.Error("invalid number of columns in first line")
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

	if table.GetMaxColumn() != 4 {
		t.Errorf("Invalid max column: %d", table.GetMaxColumn())
	}

	if table.GetMaxLine() != 72 {
		t.Errorf("Invalid max line: %d", table.GetMaxLine())
	}
}


func TestCreateTable(t *testing.T) {
	data, err1 := LoadTableFromFile("test1.csv", "auto")

	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable2(t *testing.T) {
	data, err1 := LoadTableFromFile("test1.txt", "auto")
	
	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable3(t *testing.T) {
	data, err1 := LoadTableFromFile("test1.xlsx", "auto")
	
	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData1(data, t)
}

func TestCreateTable4(t *testing.T) {
	data, err1 := LoadTableFromFile("test4.tsv", "auto")
	
	if err1 != nil {
		t.Errorf("Cannot open file: %s", err1)
		return
	}
	defer data.Close()

	CheckTestData4(data, t)
}
