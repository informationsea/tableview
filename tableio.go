package main

import "encoding/csv"

type Table interface {
	GetMaxColumn() int
	GetMaxLine() int
	GetRow(line int) []string
}

type SimpleTable struct {
	Data [][]string
	MaxColumn int
	MaxLine int
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
