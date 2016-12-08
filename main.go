package main

import ("fmt";
	"os";
	"flag";
	"github.com/nsf/termbox-go";
	"regexp";
	"errors"
)

func main() {
	var format = flag.String("format", "auto", "input format (auto/csv/tsv/tdf)")
	var fixHeader = flag.Bool("header", false, "Fix header line")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("tableview [-format FORMAT] FILE\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err2 := LoadTableFromFile(flag.Args()[0], *format)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Cannot load table file: %s\n", err2.Error())
		os.Exit(1)
	}
	defer data.Close()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	display := CreateDisplay(data, *fixHeader)
	display.Display()
	display.WaitEvent()
}

type Display struct {
	data Table
	searchText *regexp.Regexp
	voffset int
	hoffset int
	fixHeader bool
	searchMatchedLine []int
	currentMatchedLine int
}

func CreateDisplay(data Table, fixHeader bool) *Display {
	return &Display{data, nil, 0, 0, fixHeader, make([]int, 0), 0}
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
				event.Key == termbox.KeyPgdn {
				_, termHeight := termbox.Size()
				d.Scroll(termHeight, 0)
			} else if event.Ch == rune('k') || event.Ch == rune('p') ||
				event.Key == termbox.KeyCtrlP || event.Key == termbox.KeyArrowUp  {
				d.Scroll(-1, 0)
			} else if event.Ch == rune('b') || event.Ch == rune('B') ||
				event.Key == termbox.KeyCtrlB  || event.Key == termbox.KeyPgup {
				_, termHeight := termbox.Size()
				d.Scroll(-termHeight, 0)
			} else if event.Ch == rune('l') || event.Key == termbox.KeyArrowRight {
				d.Scroll(0, 1)
			} else if event.Ch == rune('h') || event.Key == termbox.KeyArrowLeft {
				d.Scroll(0, -1)
			} else if event.Ch == rune('g') || event.Key == termbox.KeyHome {
				d.ScrollTo(0, 0)
			} else if event.Ch == rune('G') || event.Key == termbox.KeyEnd {
				_, termHeight := termbox.Size()
				d.voffset = d.data.GetMaxLine() - termHeight + 1
				d.Display()
			} else if event.Ch == rune('/') {
				d.ReadSearchText()
				d.Display()
			} else if event.Ch == rune('n') {
				d.jumpSearchResultNext()
			} else if event.Ch == rune('N') {
				d.jumpSearchResultPrevious()
			}
		} else if event.Type == termbox.EventResize {
			d.Display()
		}
	}
}

func (d *Display) jumpSearchResultNext() {
	for i, v := range d.searchMatchedLine {
		if d.fixHeader {
			if d.voffset < v - 1 {
				d.currentMatchedLine = i
				d.ScrollTo(v - 1, 0);
				return
			}
		} else {
			if d.voffset < v {
				d.currentMatchedLine = i
				d.ScrollTo(v, 0);
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
	lastMatched := d.data.GetMaxLine()
	for i, v := range d.searchMatchedLine {
		if d.fixHeader {
			if d.voffset > v - 1 {
				lastMatched = i
			}
		} else {
			if d.voffset > v {
				lastMatched = i
			}
		}
	}

	if lastMatched != d.data.GetMaxLine() {
		d.currentMatchedLine = lastMatched
		if d.fixHeader {
			d.ScrollTo(d.searchMatchedLine[lastMatched] - 1, 0) 
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
	d.ScrollTo(d.voffset + v, d.hoffset + h)
}

func (d *Display) ScrollTo(newVoffset int, newHoffset int) {
	if newVoffset < 0 {
		newVoffset = 0
	} else if newVoffset >= d.data.GetMaxLine() {
		newVoffset = d.data.GetMaxLine() - 1
	}

	if newHoffset < 0 {
		newHoffset = 0
	} else if newHoffset >= d.data.GetMaxColumn() {
		newHoffset = d.data.GetMaxColumn() - 1
	}

	if newVoffset != d.voffset || newHoffset != d.hoffset {
		d.voffset = newVoffset
		d.hoffset = newHoffset
		d.Display()
	}
}

func (d *Display) ShowError(e string) {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1
	for i := 0; i < termWidth; i++ {
		termbox.SetCell(i, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for i, v := range(e) {
		termbox.SetCell(i, termHeight, v, termbox.ColorRed, termbox.ColorWhite)
	}
	termbox.SetCursor(len(e), termHeight)
	termbox.Flush()
}

func (d *Display) GetCommand(prompt string) (string, error) {
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
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
				d.searchText = nil
				return "", errors.New("Canceled")
			} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
				if len(text) > 0 {
					currentPosition -= displayWidthChar(rune(text[len(text)-1]))
					termbox.SetCell(currentPosition, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
					termbox.SetCursor(currentPosition, termHeight)
					text = text[0:len(text)-1]
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
	for i := 0; i < d.data.GetMaxLine(); i++ {
		for _, c := range d.data.GetRow(i) {
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
	_, termHeight := termbox.Size()

	lastLine := d.voffset + termHeight - 1
	firstLine := d.voffset
	if d.fixHeader {firstLine += 1}
	if lastLine > d.data.GetMaxLine() {lastLine = d.data.GetMaxLine()}

	showData := make([][]string, 0, lastLine - d.voffset)
	if (d.fixHeader) {
		showData = append(showData, d.data.GetRow(0)[d.hoffset:])
	}
	for i := firstLine; i < lastLine; i++ {
		showData = append(showData, d.data.GetRow(i)[d.hoffset:])
	}

	return showData
}

var NUMBER_RE = regexp.MustCompile("^\\d+$")

func (d *Display) Display() {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lastLine := d.voffset + termHeight
	if lastLine > d.data.GetMaxLine() {lastLine = d.data.GetMaxLine()}

	showData := d.GetDisplayData()
	
	columnSize := make([]int, d.data.GetMaxColumn())

	for _, v := range showData {
		for j := 0; j < len(v); j++ {
			textwidth := displayWidth(v[j])
			
			if (columnSize[j] < textwidth) {
				columnSize[j] = textwidth
			}
		}
	}

	termbox.SetCursor(0, 0)

	i1 := 0
	var v1 []string
	for i1, v1 = range showData {
		currentPos := 0
		
		for i2, v2 := range v1{
			if i2 != 0 {
				termbox.SetCell(currentPos + 1, i1, '|', termbox.ColorGreen, termbox.ColorDefault)
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
						termbox.ColorRed | termbox.AttrReverse | termbox.AttrBold, termbox.ColorWhite)
				} else {
					termbox.SetCell(currentPos, i1, v3, termbox.ColorDefault, termbox.ColorDefault)
				}
				
				currentPos += displayWidthChar(v3)
				if termWidth < currentPos {
					break
				}
			}

			if ! NUMBER_RE.MatchString(v2) {
				currentPos += columnSize[i2] - width
			}
		}
	}

	searchStatus := ""
	if d.searchText != nil {
		searchStatus = fmt.Sprintf("  Search:%s  Found:%d/%d", d.searchText.String(),
			d.currentMatchedLine+1, len(d.searchMatchedLine))
	}
	
	status := fmt.Sprintf("(line: %d/%d   column: %d%s)", d.voffset + 1, d.data.GetMaxLine(),
		d.hoffset + 1, searchStatus)
	currentPos := 0
	for _, v := range(status) {
		termbox.SetCell(currentPos, termHeight, v, termbox.ColorGreen, termbox.ColorDefault)
		currentPos += displayWidthChar(v)
	}
	termbox.SetCursor(currentPos, termHeight)
	termbox.Flush()
}
