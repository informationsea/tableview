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
	"testing"
)

func TestSubstringByDisplayWidth(t *testing.T) {
	actual1 := substringByDisplayWidth("hogehoge!!", 8)
	if "hogehoge" != actual1 {
		t.Errorf("Failed in test 1 / (%s)", actual1)
	}

	actual2 := substringByDisplayWidth("1日本語のヘッダ", 6)
	if "1日本" != actual2 {
		t.Errorf("Failed in test 2 / (%s)", actual2)
	}
}

func TestDisplayWidth(t *testing.T) {
	actual1 := displayWidth("hogehoge")
	if 8 != actual1 {
		t.Errorf("Failed in test 1 / (%d)", actual1)
	}

	actual2 := displayWidth("日本語のヘッダ")
	if 14 != actual2 {
		t.Errorf("Failed in test 2 / (%d)", actual2)
	}

	actual3 := displayWidth("日本語のヘッダ! ")
	if 16 != actual3 {
		t.Errorf("Failed in test 3 / (%d)", actual3)
	}
}

func TestDisplayWidthChar(t *testing.T) {
	if 1 != displayWidthChar('a') {
		t.Errorf("Failed in test 1")
	}

	if 1 != displayWidthChar(' ') {
		t.Errorf("Failed in test 2")
	}

	if 1 != displayWidthChar('1') {
		t.Errorf("Failed in test 3")
	}

	if 1 != displayWidthChar('!') {
		t.Errorf("Failed in test 4")
	}

	if 2 != displayWidthChar('日') {
		t.Errorf("Failed in test 5")
	}

	if 2 != displayWidthChar('！') {
		t.Errorf("Failed in test 6")
	}
}
