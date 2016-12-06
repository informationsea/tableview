package main

import ("fmt";
	"os";
	"flag";
	"github.com/nsf/termbox-go"
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
		panic(err2)
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
	searchText string
	voffset int
	hoffset int
	fixHeader bool
}

func CreateDisplay(data Table, fixHeader bool) *Display {
	return &Display{data, "", 0, 0, fixHeader}
}

func (d *Display) WaitEvent() {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == rune('q') || event.Key == termbox.KeyEsc {
				break
			} else if event.Ch == rune('j') || event.Ch == rune('n') ||
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
				d.hoffset = 0
				d.voffset = 0
				d.Display()
			} else if event.Ch == rune('G') || event.Key == termbox.KeyEnd {
				_, termHeight := termbox.Size()
				d.voffset = d.data.GetMaxLine() - termHeight + 1
				d.Display()
			} else if event.Ch == rune('/') {
				d.ReadSearchText()
				d.Display()
			}
		} else if event.Type == termbox.EventResize {
			d.Display()
		}
	}
}

func (d *Display) Scroll(v int, h int) {
	newVoffset := d.voffset + v
	newHoffset := d.hoffset + h

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

func (d *Display) ReadSearchText() bool {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1
	for i := 0; i < termWidth; i++ {
		termbox.SetCell(i, termHeight, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(0, termHeight, '/', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCursor(1, termHeight)

	currentPosition := 1
	text := ""
	for {
		termbox.Flush()
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
				d.searchText = ""
				return false
			} else if event.Key == termbox.KeyEnter || event.Key == termbox.KeyCtrlJ ||
				event.Key == termbox.KeyCtrlM {
				d.searchText = text
				return text != ""
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

func (d *Display) Display() {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lastLine := d.voffset + termHeight
	if lastLine > d.data.GetMaxLine() {lastLine = d.data.GetMaxLine()}

	showData := make([][]string, lastLine - d.voffset)

	if (d.fixHeader) {
		showData[0] = d.data.GetRow(0) [d.hoffset:]
		for i := d.voffset; i < lastLine-1; i++ {
			showData[i - d.voffset + 1] = d.data.GetRow(i + 1)[d.hoffset:]
		}
	} else {
		for i := d.voffset; i < lastLine; i++ {
			showData[i - d.voffset] = d.data.GetRow(i)[d.hoffset:]
		}
	}
	
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
		printLine := ""
		currentPos := 0
		
		for i2, v2 := range v1{
			if i2 != 0 {
				termbox.SetCell(currentPos + 1, i1, '|', termbox.ColorGreen, termbox.ColorDefault)
				currentPos += 3
			}

			width := displayWidth(v2)
			currentPos += columnSize[i2] - width
			
			printLine += string(v2)
			for _, v3 := range v2 {
				termbox.SetCell(currentPos, i1, v3, termbox.ColorDefault, termbox.ColorDefault)
				currentPos += displayWidthChar(v3)
				if termWidth < currentPos {
					break
				}
			}
		}
	}

	
	status := fmt.Sprintf("(line: %d/%d   column: %d)", d.voffset + 1, d.data.GetMaxLine(), d.hoffset + 1)
	currentPos := 0
	for _, v := range(status) {
		termbox.SetCell(currentPos, termHeight, v, termbox.ColorGreen, termbox.ColorDefault)
		currentPos += displayWidthChar(v)
	}
	termbox.SetCursor(currentPos, termHeight)
	termbox.Flush()
}
