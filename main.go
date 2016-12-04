package main

import ("fmt";
	"os";
	"encoding/csv";
	"io";
	"golang.org/x/text/width"
)

func displayWidth(text string) int {
	textwidth := 0
	for _, value := range text {
		p, _ := width.LookupString(string([]rune{value}))
		if p.Kind() == width.EastAsianAmbiguous || p.Kind() == width.EastAsianWide || p.Kind() == width.EastAsianFullwidth {
			textwidth += 2
		} else {
			textwidth += 1
		}
	}
	return textwidth
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("tableview FILE\n")
		os.Exit(1)
	}

	inputFile, _ := os.Open(os.Args[1])
	csvReader := csv.NewReader(inputFile)

	data, _ := csvReader.ReadAll()

	columnSize := make([]int, len(data[0]))

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			textwidth := displayWidth(data[i][j])
			
			if (columnSize[j] < textwidth) {
				columnSize[j] = textwidth
			}
		}
	}

	for _, v1 := range data {
		for i2, v2 := range v1 {
			width := len(v2)
			if i2 != 0 {
				fmt.Print(" | ")
			}

			for j := 0; j < columnSize[i2] - width; j++ {
				fmt.Print(" ")
			}
			
			fmt.Printf("%s", v2)
		}
		fmt.Println()
	}
}

type TableView struct {
	file io.Reader
	csvReader csv.Reader
}
