package main

import (
	"fmt"
	"os"
	"strings"
)

// Banner is a map that holds ASCII art for each character
// each character maps to a slice of 8 strings (the art lines)
type Banner map[rune][]string

const (
	firstChar  = 32  // space character (first printable ASCII)
	lastChar   = 126 // ~ character (last printable ASCII)
	charHeight = 8   // each character is 8 lines tall
)

// LoadBanner reads a banner file and loads all the character art
// banner files have each character as 9 lines (1 empty + 8 art)
func LoadBanner(path string) (Banner, error) {
	// read the whole file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Convert to string and handle both Unix (\n) and Windows (\r\n) line endings
	// Replace \r\n with \n first, then remove any remaining \r
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	// split file into lines
	lines := strings.Split(content, "\n")

	// create the banner map
	banner := make(Banner)

	const blockSize = charHeight + 1 // each char = 9 lines total

	// go through each ASCII character from space to ~
	for code := firstChar; code <= lastChar; code++ {
		// figure out where this character starts in the file
		blockIndex := int(code - firstChar)
		start := blockIndex * blockSize

		// make sure we have enough lines in the file
		if start+1+charHeight > len(lines) {
			return nil, fmt.Errorf("invalid banner file: not enough lines for char %q", rune(code))
		}

		// skip the first empty line, take the next 8 lines
		glyphLines := lines[start+1 : start+1+charHeight]
		banner[rune(code)] = glyphLines
	}

	return banner, nil
}
