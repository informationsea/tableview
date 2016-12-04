package main

import ("fmt";
	"os";
	"encoding/csv";
	"golang.org/x/text/width";
	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("tableview FILE\n")
		os.Exit(1)
	}

	inputFile, err3 := os.Open(os.Args[1])
	if err3 != nil {
		panic(err3)
	}
	
	defer inputFile.Close()
	csvReader := csv.NewReader(inputFile)

	data, err2 := csvReader.ReadAll()
	if err2 != nil {
		panic(err2)
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	maxColumn := 0
	for _, v1 := range data {
		if maxColumn < len(v1) {
			maxColumn = len(v1)
		}
	}

	voffset := 0
	hoffset := 0
	display(data, voffset, hoffset)

	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == rune('q') || event.Key == termbox.KeyEsc {
				break
			} else if event.Ch == rune('j') || event.Ch == rune('n') || event.Key == termbox.KeyCtrlN || event.Key == termbox.KeyArrowDown || event.Key == termbox.KeyEnter {
				voffset += 1
				if voffset >= len(data) {voffset = len(data)-1}
				display(data, voffset, hoffset)
			} else if event.Ch == rune('F') || event.Ch == rune('f') || event.Key == termbox.KeyCtrlV || event.Key == termbox.KeyCtrlF || event.Key == termbox.KeyPgdn {
				_, termHeight := termbox.Size()
				voffset += termHeight
				if voffset >= len(data) {voffset = len(data)-1}
				display(data, voffset, hoffset)
			} else if event.Ch == rune('k') || event.Ch == rune('p') || event.Key == termbox.KeyCtrlP || event.Key == termbox.KeyArrowUp  {
				voffset -= 1
				if voffset < 0 {voffset = 0}
				display(data, voffset, hoffset)
			} else if event.Ch == rune('b') || event.Ch == rune('B') || event.Key == termbox.KeyCtrlB  || event.Key == termbox.KeyPgup {
				_, termHeight := termbox.Size()
				voffset -= termHeight
				if voffset < 0 {voffset = 0}
				display(data, voffset, hoffset)
			} else if event.Ch == rune('l') || event.Key == termbox.KeyArrowRight {
				hoffset += 1
				if hoffset >= maxColumn {hoffset = maxColumn - 1}
				display(data, voffset, hoffset)
			} else if event.Ch == rune('h') || event.Key == termbox.KeyArrowLeft {
				hoffset -= 1
				if hoffset < 0 {hoffset = 0}
				display(data, voffset, hoffset)
			}
			
			
		} else if event.Type == termbox.EventResize {
			display(data, voffset, hoffset)
		}
	}
}

func display(data [][]string, offset int, hoffset int) {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	columnSize := make([]int, len(data[0]))

	for i := offset; i < len(data) && (i-offset) < termHeight; i++ {
		for j := 0; j < len(data[i]); j++ {
			textwidth := displayWidth(data[i][j])
			
			if (columnSize[j] < textwidth) {
				columnSize[j] = textwidth
			}
		}
	}

	termbox.SetCursor(0, 0)

	i1 := 0
	for i1 = offset; i1 < len(data) && (i1-offset) < termHeight; i1++ {
		printLine := ""
		for i2 := hoffset; i2 < len(data[i1]); i2++ {
			v2 := data[i1][i2]
			width := displayWidth(v2)
			if i2 != hoffset {
				printLine += " | "
			}

			for j := 0; j < columnSize[i2] - width; j++ {
				printLine += " "
			}
			
			printLine += v2
		}

		fmt.Print(substringByDisplayWidth(printLine, termWidth))

		for i2 := displayWidth(printLine); i2 < termWidth; i2++ {
			fmt.Print(" ")
		}
		
		fmt.Println()
	}

	for ; (i1-offset) < termHeight; i1++ {
		for i2 := 0; i2 < termWidth; i2++ {
			fmt.Print(" ")
		}
		fmt.Println()
	}

	fmt.Printf("(line: %d / column: %d)   ", offset, hoffset)
	termbox.SetCursor(0, termHeight)
	termbox.Flush()
}

func substringByDisplayWidth(text string, width int) string {
	newText := ""
	textWidth := 0
	for _, value := range text {
		textWidth += displayWidthChar(rune(value))
		if textWidth < width {
			newText += string(value)
		} else {
			break
		}
	}
	return newText
}


func displayWidth(text string) int {
	textwidth := 0
	for _, value := range text {
		textwidth += displayWidthChar(rune(value))
	}
	return textwidth
}

func displayWidthChar(char rune) int {
	textwidth := 0
	p, _ := width.LookupString(string([]rune{char}))
	if p.Kind() == width.EastAsianAmbiguous || p.Kind() == width.EastAsianWide || p.Kind() == width.EastAsianFullwidth {
		textwidth += 2
	} else {
		textwidth += 1
	}
	return textwidth
}
