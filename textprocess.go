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
