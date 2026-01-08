package main

import (
	"fmt"
	"testing"
)

// make a fake banner for testing
func fakeBanner() Banner {
	b := make(Banner)

	// create test data for A, B and space
	for i := 0; i < charHeight; i++ {
		if i == 0 {
			// initialize the slices first
			b['A'] = make([]string, charHeight)
			b['B'] = make([]string, charHeight)
			b[' '] = make([]string, charHeight)
		}
		// fill each line with test data
		b['A'][i] = fmt.Sprintf("A%d", i)
		b['B'][i] = fmt.Sprintf("B%d", i)
		b[' '][i] = fmt.Sprintf(" %d", i)
	}

	return b
}

// test loading the standard banner file
func TestLoadBannerStandard(t *testing.T) {
	banner, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("failed to load banner: %v", err)
	}

	// make sure we got something
	if len(banner) == 0 {
		t.Fatalf("banner is empty")
	}

	// should have 95 characters (32-126)
	if len(banner) != 95 {
		t.Errorf("expected 95 characters, got %d", len(banner))
	}

	// check each character exists and has right height
	for c := firstChar; c <= lastChar; c++ {
		art, ok := banner[rune(c)]
		if !ok {
			t.Errorf("character %q not found", c)
			continue
		}
		// each character should be 8 lines tall
		if len(art) != charHeight {
			t.Errorf("character %q has wrong height: %d", c, len(art))
		}
	}
}

// test rendering a simple line
func TestRenderLine(t *testing.T) {
	b := fakeBanner()

	// render "AB" and check output
	got := RenderLine("AB", b)
	// expected output: A and B side by side, 8 rows
	want := "" +
		"A0B0\n" +
		"A1B1\n" +
		"A2B2\n" +
		"A3B3\n" +
		"A4B4\n" +
		"A5B5\n" +
		"A6B6\n" +
		"A7B7\n"

	if got != want {
		t.Errorf("RenderLine(\"AB\") =\n%q\nwant:\n%q", got, want)
	}
}

// test with character not in banner
func TestRenderLineUnknownChar(t *testing.T) {
	b := fakeBanner()

	// Z is not in our fake banner, should not crash
	_ = RenderLine("Z", b)
}

// test converting \n to real newlines
func TestDecodeEscapedNewlines(t *testing.T) {
	// input has literal \n, should become real newlines
	got := decodeEscapedNewlines("hello\\nworld\\n")
	want := "hello\nworld\n"

	if got != want {
		t.Errorf("decodeEscapedNewlines(...) = %q, want %q", got, want)
	}
}

// test rendering input with space
func TestRenderInputSingleLine(t *testing.T) {
	b := fakeBanner()

	// "A B" should render A, space, B
	got := RenderInput("A B", b)
	// A + space + B on each row
	want := "" +
		"A0 0B0\n" +
		"A1 1B1\n" +
		"A2 2B2\n" +
		"A3 3B3\n" +
		"A4 4B4\n" +
		"A5 5B5\n" +
		"A6 6B6\n" +
		"A7 7B7\n"

	if got != want {
		t.Errorf("RenderInput(\"A B\") =\n%q\nwant:\n%q", got, want)
	}
}

// test with newlines in input
func TestRenderInputWithLogicalNewlines(t *testing.T) {
	b := fakeBanner()

	// "A\nB" should render A on first line, B on second
	got := RenderInput("A\\nB", b)
	// A rendered completely, then B rendered completely
	want := "" +
		"A0\n" +
		"A1\n" +
		"A2\n" +
		"A3\n" +
		"A4\n" +
		"A5\n" +
		"A6\n" +
		"A7\n" +
		"B0\n" +
		"B1\n" +
		"B2\n" +
		"B3\n" +
		"B4\n" +
		"B5\n" +
		"B6\n" +
		"B7\n"

	if got != want {
		t.Errorf("RenderInput(\"A\\nB\") =\n%q\nwant:\n%q", got, want)
	}
}

// test empty input
func TestRenderInputEmpty(t *testing.T) {
	b := fakeBanner()

	// empty input should give empty output
	got := RenderInput("", b)
	if got != "" {
		t.Errorf("RenderInput(\"\") = %q, want empty string", got)
	}
}
