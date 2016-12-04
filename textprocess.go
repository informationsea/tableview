package main

import (
	"golang.org/x/text/width";
)

func substringByDisplayWidth(text string, width int) string {
	newText := ""
	textWidth := 0
	for _, value := range text {
		textWidth += displayWidthChar(rune(value))
		if textWidth <= width {
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
