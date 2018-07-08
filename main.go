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
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/nsf/termbox-go"
)

const VERSION = "@DEV@"

func ShowLicense() {
	license, _ := Asset("LICENSE")
	reader := strings.NewReader(string(license))
	data := CreatePartialTable(reader, ParseTSVRecord)
	display := CreateDisplay(data, false)
	display.helpMode = true
	display.loadAllData()
	display.Display()
	display.WaitEvent()
}

var logger *log.Logger = log.New(EmptyWriter{}, "", log.LstdFlags)

func main() {
	var format = flag.String("format", "auto", "input format (auto/csv/tsv/tdf)")
	var sheetNum = flag.Int("sheet", 1, "Sheet index (Excel only)")
	var fixHeader = flag.Bool("header", false, "Fix header line")
	var showVersion = flag.Bool("version", false, "Show version")
	var showLicense = flag.Bool("license", false, "Show license")
	var showHelp = flag.Bool("help", false, "Show help")
	var loggerPath = flag.String("logger", "", "Log file for debug")
	flag.Parse()

	if *loggerPath != "" {
		logFile, err := os.OpenFile(*loggerPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		defer logFile.Close()

		logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if *showLicense {
		err := termbox.Init()
		if err != nil {
			panic(err)
		}
		defer termbox.Close()
		ShowLicense()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("tableview : human friendly table viewer\nVersion: %s\n\n", VERSION)
		fmt.Println("Copyright (C) 2016  OKAMURA, Yasunobu")
		fmt.Println("")
		fmt.Println("This program is free software: you can redistribute it and/or modify")
		fmt.Println("it under the terms of the GNU General Public License as published by")
		fmt.Println("the Free Software Foundation, either version 3 of the License, or")
		fmt.Println("(at your option) any later version.")
		fmt.Println("")
		fmt.Println("This program is distributed in the hope that it will be useful,")
		fmt.Println("but WITHOUT ANY WARRANTY; without even the implied warranty of")
		fmt.Println("MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the")
		fmt.Println("GNU General Public License for more details.")
		fmt.Println("")
		fmt.Println("You should have received a copy of the GNU General Public License")
		fmt.Println("along with this program.  If not, see <http://www.gnu.org/licenses/>.")
		os.Exit(0)
	}

	if *showHelp || len(flag.Args()) > 1 || (len(flag.Args()) == 0 && isatty.IsTerminal(os.Stdin.Fd())) {
		fmt.Println("tableview [-format FORMAT] [FILE]")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()

		if *showHelp {
			for _, v := range HELP[4:] {
				fmt.Println(v[0])
			}
			os.Exit(0)
		} else {
			os.Exit(1)
		}

	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(os.Stderr, "No output terminal is found\n")
		os.Exit(1)
	}

	//fmt.Println("Now loading...")
	var data Table
	var err error

	if len(flag.Args()) == 0 {
		if *format == "csv" {
			csvReader := csv.NewReader(os.Stdin)
			data = CreatePartialCSV(csvReader)
		} else {
			data = CreatePartialTable(os.Stdin, ParseTSVRecord)
		}
	} else {
		data, err = LoadTableFromFile(flag.Args()[0], *format, *sheetNum)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load table file: %s\n", err.Error())
		os.Exit(1)
	}

	defer data.Close()

	_, err = data.GetRow(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Empty data\n")
		return
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	display := CreateDisplay(data, *fixHeader)
	display.Display()
	display.WaitEvent()
}

var HELP = [][]string{[]string{"tableview @DEV@"},
	[]string{"Copyright (C) 2016 OKAMURA, Yasunobu"},
	[]string{""},
	[]string{""},
	[]string{"Key Binding"},
	[]string{""},
	[]string{"  j, Ctrl-N, ARROW-DOWN, ENTER"},
	[]string{"       Scroll forward 1 line"},
	[]string{""},
	[]string{"  k, Ctrl-P, ARROW-UP"},
	[]string{"       Scroll back 1 line"},
	[]string{""},
	[]string{"  h, ARROW-LEFT"},
	[]string{"       Scroll horizontally left 1 column"},
	[]string{""},
	[]string{"  l, ARROW-RIGHT"},
	[]string{"       Scroll horizontally right 1 column"},
	[]string{""},
	[]string{"  f, Ctrl-V, SPACE, PageDown"},
	[]string{"       Scroll forward 1 window"},
	[]string{""},
	[]string{"  ["},
	[]string{"       Scroll horizontally left 1 character"},
	[]string{""},
	[]string{"  ]"},
	[]string{"       Scroll horizontally right 1 character"},
	[]string{""},
	[]string{"  g, Home, <"},
	[]string{"       Go to first line"},
	[]string{""},
	[]string{"  G, End, >"},
	[]string{"       Go to last line"},
	[]string{""},
	[]string{"  /pattern"},
	[]string{"       Search forward in the file"},
	[]string{""},
	[]string{"  n"},
	[]string{"       Repeat previous search"},
	[]string{""},
	[]string{"  N"},
	[]string{"       Repeat previous search, but in the reverse direction"},
	[]string{""},
	[]string{"  ?"},
	[]string{"       Show this help"},
	[]string{""},
	[]string{"  L"},
	[]string{"       Show license"},
}

type Display struct {
	data               Table
	searchText         *regexp.Regexp
	voffset            int
	hoffset            int
	hoffset2           int
	fixHeader          bool
	searchMatchedLine  []int
	currentMatchedLine int
	helpMode           bool
}

func CreateDisplay(data Table, fixHeader bool) *Display {
	return &Display{data, nil, 0, 0, 0, fixHeader, make([]int, 0), 0, false}
}

func (d *Display) WaitEvent() {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == rune('q') || event.Key == termbox.KeyEsc {
				break
			} else if event.Ch == rune('j') ||
				event.Key == termbox.KeyCtrlN || event.Key == termbox.KeyArrowDown ||
				event.Key == termbox.KeyEnter {
				d.Scroll(1, 0)
			} else if event.Ch == rune('F') || event.Ch == rune('f') ||
				event.Key == termbox.KeyCtrlV || event.Key == termbox.KeyCtrlF ||
				event.Key == termbox.KeyPgdn || event.Key == termbox.KeySpace {
				_, termHeight := termbox.Size()
				d.Scroll(termHeight, 0)
			} else if event.Ch == rune('k') || event.Ch == rune('p') ||
				event.Key == termbox.KeyCtrlP || event.Key == termbox.KeyArrowUp {
				d.Scroll(-1, 0)
			} else if event.Ch == rune('b') || event.Ch == rune('B') ||
				event.Key == termbox.KeyCtrlB || event.Key == termbox.KeyPgup {
				_, termHeight := termbox.Size()
				d.Scroll(-termHeight, 0)
			} else if event.Ch == rune('l') || event.Key == termbox.KeyArrowRight {
				d.hoffset2 = 0
				d.Scroll(0, 1)
			} else if event.Ch == rune('h') || event.Key == termbox.KeyArrowLeft {
				d.hoffset2 = 0
				d.Scroll(0, -1)
			} else if event.Ch == rune('[') {
				d.hoffset2 -= 1
				if d.hoffset2 < 0 {
					d.hoffset2 = 0
				}
				d.Display()
			} else if event.Ch == rune(']') {
				d.hoffset2 += 1
				d.Display()
			} else if event.Ch == rune('g') || event.Key == termbox.KeyHome || event.Ch == rune('<') {
				d.hoffset2 = 0
				d.ScrollTo(0, 0)
			} else if event.Ch == rune('G') || event.Key == termbox.KeyEnd || event.Ch == rune('>') {
				_, termHeight := termbox.Size()
				err := d.loadAllData()
				if err != nil {
					d.ShowStatus("Canceled by user", termbox.ColorBlue)
				} else {
					count, err2 := d.data.GetLineCountIfAvailable()
					if err2 != nil {
						panic("Cannot get line count")
					}
					d.voffset = count - termHeight + 1
					d.Display()
				}
			} else if event.Ch == rune('/') {
				d.ReadSearchText()
				d.Display()
			} else if event.Ch == rune('n') {
				d.jumpSearchResultNext()
			} else if event.Ch == rune('N') {
				d.jumpSearchResultPrevious()
			} else if event.Ch == rune('?') {
				if !d.helpMode {
					display := CreateDisplay(CreateTable(HELP), false)
					display.helpMode = true
					display.Display()
					display.WaitEvent()
					d.Display()
				}
			} else if event.Ch == rune('L') {
				ShowLicense()
				d.Display()
			}
		} else if event.Type == termbox.EventResize {
			d.Display()
		}
	}
}

func (d *Display) loadAllData() error {
	_, err := d.data.GetLineCountIfAvailable()
	if err == nil {
		return nil
	}

	d.ShowStatus("Now loading...", termbox.ColorBlue)
	finish := make(chan bool)
	cancel := make(chan bool)

	go func() {
		for {
			c, e := d.data.LoadAll(1000)
			if e != nil {
				panic(e)
			}
			if !c {
				break
			}

			select {
			case <-cancel:
				finish <- false
				return
			default:
				// do nothing
			}
		}
		finish <- true
	}()

	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventKey {
				if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
					break
				}
			}
		}
		cancel <- true
	}()

	status := <-finish
	if status {
		termbox.Interrupt()
		return nil
	}
	return errors.New("Interrupt")
}

func (d *Display) jumpSearchResultNext() {
	for i, v := range d.searchMatchedLine {
		if d.fixHeader {
			if d.voffset < v-1 {
				d.currentMatchedLine = i
				d.ScrollTo(v-1, 0)
				return
			}
		} else {
			if d.voffset < v {
				d.currentMatchedLine = i
				d.ScrollTo(v, 0)
				return
			}
		}
	}

	if len(d.searchMatchedLine) > 0 {
		d.ShowError("No more found line")
	} else {
		d.ShowError("No found line")
	}
}

func (d *Display) jumpSearchResultPrevious() {
	lastMatched := d.data.GetLoadedLineCount()
	for i, v := range d.searchMatchedLine {
		if d.fixHeader {
			if d.voffset > v-1 {
				lastMatched = i
			}
		} else {
			if d.voffset > v {
				lastMatched = i
			}
		}
	}

	if lastMatched != d.data.GetLoadedLineCount() {
		d.currentMatchedLine = lastMatched
		if d.fixHeader {
			d.ScrollTo(d.searchMatchedLine[lastMatched]-1, 0)
		} else {
			d.ScrollTo(d.searchMatchedLine[lastMatched], 0)
		}
		return
	}

	if len(d.searchMatchedLine) > 0 {
		d.ShowError("No more found line")
	} else {
		d.ShowError("No found line")
	}
}

func (d *Display) Scroll(v int, h int) {
	d.ScrollTo(d.voffset+v, d.hoffset+h)
}

func (d *Display) ScrollTo(newVoffset int, newHoffset int) {
	if newVoffset < 0 {
		newVoffset = 0
	} else {
		count, err := d.data.GetLineCountIfAvailable()
		if err == nil && newVoffset >= count {
			newVoffset = count - 1
		}
		if newVoffset < 0 {
			newVoffset = 0
		}
	}

	if newHoffset < 0 {
		newHoffset = 0
	}

	if newVoffset != d.voffset || newHoffset != d.hoffset {
		d.voffset = newVoffset
		d.hoffset = newHoffset
		d.Display()
	}
}

func (d *Display) ShowError(e string) {
	d.ShowStatus(e, termbox.ColorRed)
}

func (d *Display) ShowStatus(e string, color termbox.Attribute) {
	logger.Printf("ShowStatus %s\n", e)
	termWidth, termHeight := termbox.Size()
	termHeight -= 1
	for i := 0; i < termWidth; i++ {
		termbox.SetCell(i, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for i, v := range e {
		termbox.SetCell(i, termHeight, v, color, termbox.ColorWhite)
	}
	termbox.SetCursor(len(e), termHeight)
	termbox.Flush()
}

func (d *Display) GetCommand(prompt string) (string, error) {
	logger.Println("GetCommand")
	termWidth, termHeight := termbox.Size()
	termHeight -= 1
	for i := 0; i < termWidth; i++ {
		termbox.SetCell(i, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for _, v := range prompt {
		termbox.SetCell(0, termHeight, v, termbox.ColorGreen, termbox.ColorDefault)
	}
	termbox.SetCursor(1, termHeight)

	currentPosition := len(prompt)
	text := ""
	for {
		termbox.Flush()
		event := termbox.PollEvent()
		logger.Printf("Event: %c\n", event.Ch)
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
				d.searchText = nil
				return "", errors.New("Canceled")
			} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
				if len(text) > 0 {
					currentPosition -= displayWidthChar(rune(text[len(text)-1]))
					termbox.SetCell(currentPosition, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
					termbox.SetCursor(currentPosition, termHeight)
					text = text[0 : len(text)-1]
				}
			} else if event.Key == termbox.KeyEnter || event.Key == termbox.KeyCtrlJ ||
				event.Key == termbox.KeyCtrlM {

				return text, nil

			} else {
				text += string(event.Ch)
				termbox.SetCell(currentPosition, termHeight, event.Ch,
					termbox.ColorGreen, termbox.ColorDefault)
				currentPosition += displayWidthChar(event.Ch)
				termbox.SetCursor(currentPosition, termHeight)
			}
		}
	}
}

func (d *Display) ReadSearchText() bool {
	text, err1 := d.GetCommand("/")
	if err1 != nil {
		d.searchText = nil
		d.ShowError(err1.Error())
		return false
	}

	if text == "" {
		d.searchText = nil
		return false
	}

	reg, err2 := regexp.Compile(text)

	if err2 != nil {
		d.searchText = nil
		d.ShowError(err2.Error())
		return false
	}

	d.searchText = reg
	d.searchMatchedLine = make([]int, 0)
	d.currentMatchedLine = -1

	err := d.loadAllData()
	if err != nil {
		d.ShowStatus("Canceled by user", termbox.ColorBlue)
		return false
	}

	d.ShowStatus("Now searching...", termbox.ColorBlue)
	for i := 0; i < d.data.GetLoadedLineCount(); i++ {
		row, err := d.data.GetRow(i)
		if err != nil {
			panic(err)
		}

		for _, c := range row {
			if d.searchText.FindString(c) != "" {
				d.searchMatchedLine = append(d.searchMatchedLine, i)
				break
			}
		}
	}

	d.jumpSearchResultNext()

	return true
}

func (d *Display) GetDisplayData() [][]string {
	logger.Println("GetDisplayData")
	_, termHeight := termbox.Size()

	lastLine := d.voffset + termHeight - 1
	firstLine := d.voffset
	if d.fixHeader {
		firstLine++
	}
	logger.Println("Start loading data")
	err := d.data.Load(lastLine)

	if err != nil {
		logger.Printf("err: %s\n", err.Error())
		termbox.Close()
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		logger.Println("TermBox closed")
		os.Exit(1)
	}

	if lastLine > d.data.GetLoadedLineCount() {
		lastLine = d.data.GetLoadedLineCount()
	}

	logger.Println("Data loaded")

	showData := make([][]string, 0, lastLine-d.voffset)
	if d.fixHeader {
		row, err := d.data.GetRow(0)

		if err != nil {
			showData = append(showData, []string{})
		} else if d.hoffset < len(row) {
			showData = append(showData, row[d.hoffset:])
		} else {
			showData = append(showData, []string{})
		}
		logger.Printf("header row: %s\n", row)
	}
	logger.Println("load header")
	for i := firstLine; i < lastLine; i++ {
		row, err := d.data.GetRow(i)

		if err != nil {
			showData = append(showData, []string{})
			break
		}

		if d.hoffset < len(row) {
			showData = append(showData, row[d.hoffset:])
		} else {
			showData = append(showData, []string{})
		}
		logger.Printf("row: %s\n", row)
	}
	logger.Println("done")

	return showData
}

var NUMBER_RE = regexp.MustCompile("^[\\d\\.-]+$")

func (d *Display) Display() {
	logger.Printf("Display\n")
	termWidth, termHeight := termbox.Size()
	termHeight -= 1

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	showData := d.GetDisplayData()

	columnSize := make([]int, 0)

	for _, v := range showData {
		for j := len(columnSize); j < len(v); j++ {
			columnSize = append(columnSize, 0)
		}

		for j := 0; j < len(v); j++ {
			textwidth := displayWidth(v[j])

			if columnSize[j] < textwidth {
				columnSize[j] = textwidth
			}
		}
	}

	if len(columnSize) == 0 {
		if d.hoffset != 0 {
			d.Scroll(0, -1)
			return
		}
		d.hoffset2 = 0
	} else {
		if columnSize[0] <= d.hoffset2 {
			d.hoffset2 = columnSize[0] - 1
		}
	}

	termbox.SetCursor(0, 0)

	i1 := 0
	var v1 []string
	for i1, v1 = range showData {
		currentPos := -d.hoffset2

		for i2, v2 := range v1 {
			if i2 != 0 {
				termbox.SetCell(currentPos+1, i1, '|', termbox.ColorGreen, termbox.ColorDefault)
				currentPos += 3
			}

			if NUMBER_RE.MatchString(v2) {

			}

			width := displayWidth(v2)

			if NUMBER_RE.MatchString(v2) {
				currentPos += columnSize[i2] - width
			}

			var matches [][]int
			if d.searchText != nil {
				matches = d.searchText.FindAllStringIndex(v2, -1)
			}

			for i3, v3 := range v2 {
				searchMatch := false
				for _, a := range matches {
					if a[0] <= i3 && i3 < a[1] {
						searchMatch = true
					}
				}

				if searchMatch {
					termbox.SetCell(currentPos, i1, v3,
						termbox.ColorRed|termbox.AttrReverse|termbox.AttrBold, termbox.ColorWhite)
				} else {
					termbox.SetCell(currentPos, i1, v3, termbox.ColorDefault, termbox.ColorDefault)
				}

				currentPos += displayWidthChar(v3)
				if termWidth <= currentPos+1 {
					termbox.SetCell(termWidth-1, i1, '>', termbox.ColorRed, termbox.ColorDefault)
					break
				}
			}

			if !NUMBER_RE.MatchString(v2) {
				currentPos += columnSize[i2] - width
			}
		}
	}

	searchStatus := ""
	if d.searchText != nil {
		searchStatus = fmt.Sprintf("  Search:%s  Found:%d/%d", d.searchText.String(),
			d.currentMatchedLine+1, len(d.searchMatchedLine))
	}

	count, err := d.data.GetLineCountIfAvailable()
	status := ""
	if err == nil {
		status = fmt.Sprintf("(line: %d/%d   column: %d/%d%s)", d.voffset+1, count,
			d.hoffset+1, len(columnSize)+d.hoffset, searchStatus)
	} else {
		status = fmt.Sprintf("(line: %d/%d+   column: %d/%d%s)", d.voffset+1, d.data.GetLoadedLineCount(),
			d.hoffset+1, len(columnSize)+d.hoffset, searchStatus)
	}

	currentPos := 0
	for _, v := range status {
		termbox.SetCell(currentPos, termHeight, v, termbox.ColorGreen, termbox.ColorDefault)
		currentPos += displayWidthChar(v)
	}
	termbox.SetCursor(currentPos, termHeight)
	termbox.Flush()
}
