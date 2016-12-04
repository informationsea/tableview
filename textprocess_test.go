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

