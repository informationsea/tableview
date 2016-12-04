package main

import ("fmt";
	"os";
	"encoding/csv";
	"strings";
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

	filename := flag.Args()[0]

	if (*format == "auto") {
		if strings.HasSuffix(filename, ".csv") {
			*format = "csv"
		} else if strings.HasSuffix(filename, ".txt") {
			*format = "tsv"
		} else if strings.HasSuffix(filename, ".tsv") {
			*format = "tsv"
		} else if strings.HasSuffix(filename, ".tdf") {
			*format = "tsv"
		} else {
			fmt.Println("Cannot suggest format")
			fmt.Println("Please set -format flag")
			os.Exit(1)
		}
	} else if (*format == "tdf") {
		*format = "tsv"
	} else if !(*format == "csv" || *format == "tsv") {
		fmt.Printf("Invalid format: %s\n", *format)
		flag.PrintDefaults()
		os.Exit(1)
	}


	inputFile, err3 := os.Open(flag.Args()[0])
	if err3 != nil {
		panic(err3)
	}
	
	defer inputFile.Close()
	csvReader := csv.NewReader(inputFile)
	if *format == "tsv" {
		csvReader.Comma = '\t'
		csvReader.LazyQuotes = true
	}

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
	display(data, voffset, hoffset, *fixHeader)

	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == rune('q') || event.Key == termbox.KeyEsc {
				break
			} else if event.Ch == rune('j') || event.Ch == rune('n') ||
				event.Key == termbox.KeyCtrlN || event.Key == termbox.KeyArrowDown ||
				event.Key == termbox.KeyEnter {
				voffset += 1
				if voffset >= len(data) {voffset = len(data)-1}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('F') || event.Ch == rune('f') ||
				event.Key == termbox.KeyCtrlV || event.Key == termbox.KeyCtrlF ||
				event.Key == termbox.KeyPgdn {
				_, termHeight := termbox.Size()
				voffset += termHeight
				if voffset >= len(data) {voffset = len(data)-1}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('k') || event.Ch == rune('p') ||
				event.Key == termbox.KeyCtrlP || event.Key == termbox.KeyArrowUp  {
				voffset -= 1
				if voffset < 0 {voffset = 0}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('b') || event.Ch == rune('B') ||
				event.Key == termbox.KeyCtrlB  || event.Key == termbox.KeyPgup {
				_, termHeight := termbox.Size()
				voffset -= termHeight
				if voffset < 0 {voffset = 0}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('l') || event.Key == termbox.KeyArrowRight {
				hoffset += 1
				if hoffset >= maxColumn {hoffset = maxColumn - 1}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('h') || event.Key == termbox.KeyArrowLeft {
				hoffset -= 1
				if hoffset < 0 {hoffset = 0}
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('g') || event.Key == termbox.KeyHome {
				hoffset = 0
				voffset = 0
				display(data, voffset, hoffset, *fixHeader)
			} else if event.Ch == rune('G') || event.Key == termbox.KeyEnd {
				_, termHeight := termbox.Size()
				voffset = len(data) - termHeight + 1
				display(data, voffset, hoffset, *fixHeader)
			}
		} else if event.Type == termbox.EventResize {
			display(data, voffset, hoffset, *fixHeader)
		}
	}
}


func display(data [][]string, offset int, hoffset int, fixHeader bool) {
	termWidth, termHeight := termbox.Size()
	termHeight -= 1

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lastLine := offset + termHeight
	if lastLine > len(data) {lastLine = len(data)}

	showData := make([][]string, lastLine - offset)

	if (fixHeader) {
		showData[0] = data[0][hoffset:]
		for i := offset; i < lastLine-1; i++ {
			showData[i - offset + 1] = data[i + 1][hoffset:]
		}
	} else {
		showData = data[offset:lastLine][hoffset:]
	}
	
	columnSize := make([]int, len(data[0]))

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
		for i2, v2 := range v1{
			width := displayWidth(v2)
			if i2 != hoffset {
				printLine += " | "
			}

			for j := 0; j < columnSize[i2] - width; j++ {
				printLine += " "
			}
			
			printLine += string(v2)
		}

		fmt.Print(substringByDisplayWidth(printLine, termWidth))

		for i2 := displayWidth(printLine); i2 < termWidth; i2++ {
			fmt.Print(" ")
		}
		
		fmt.Println()
	}

	for ; i1 < termHeight-1; i1++ {
		for i2 := 0; i2 < termWidth; i2++ {
			fmt.Print(" ")
		}
		fmt.Println()
	}

	fmt.Printf("(line: %d/%d   column: %d)   ", offset + 1, len(data), hoffset + 1)
	termbox.SetCursor(0, termHeight)
	termbox.Flush()
}

