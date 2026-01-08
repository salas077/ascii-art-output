package main

import "strings"

// RenderLine takes a string and makes ASCII art from it
// builds it row by row (each character is 8 rows tall)
func RenderLine(s string, b Banner) string {
	const height = charHeight

	var builder strings.Builder

	// make empty glyph for characters we don't have
	empty := make([]string, height)

	// go through each row (0 to 7)
	for row := 0; row < height; row++ {
		// for each character in the input string
		for _, ch := range s {
			glyph, ok := b[ch]
			if !ok {
				// character not in banner, use empty space
				glyph = empty
			}
			// add this character's row to the line
			builder.WriteString(glyph[row])
		}
		// end of row, add newline
		builder.WriteRune('\n')
	}

	return builder.String()
}

// decodeEscapedNewlines converts \n to actual newlines
// when user types "Hello\nWorld" we want real newlines
func decodeEscapedNewlines(s string) string {
	runes := []rune(s)
	out := make([]rune, 0, len(runes))

	// go through each character
	for i := 0; i < len(runes); i++ {
		// if we see \ followed by n, make it a real newline
		if runes[i] == '\\' && i+1 < len(runes) && runes[i+1] == 'n' {
			out = append(out, '\n')
			i++ // skip the n since we already used it
		} else {
			// normal character, just copy it
			out = append(out, runes[i])
		}
	}

	return string(out)
}

// RenderInput is the main function that handles everything
// takes user input and converts it to ASCII art
func RenderInput(input string, b Banner) string {
	// first convert \n strings to real newlines
	input = decodeEscapedNewlines(input)

	// if input is empty, return empty
	if input == "" {
		return ""
	}

	// if no newlines, just render as single line
	if !strings.Contains(input, "\n") {
		return RenderLine(input, b)
	}

	var builder strings.Builder

	// split input by newlines to handle multiple lines
	parts := strings.Split(input, "\n")
	hadText := false

	// process each line
	for i, part := range parts {
		last := i == len(parts)-1

		if part != "" {
			// non-empty line, render it
			asciiBlock := RenderLine(part, b)
			builder.WriteString(asciiBlock)
			hadText = true

		} else {
			// empty line handling
			if !last {
				// empty line in middle, add blank line
				builder.WriteRune('\n')
			} else if hadText {
				// empty line at end, keep the newline
				builder.WriteRune('\n')
			}
		}
	}

	return builder.String()
}

// RenderWithColor renders text with specific characters colored
// It works character-by-character, adding color codes where needed
//
// Parameters:
// - input: the text to render (e.g., "a kitten")
// - banner: the loaded banner map
// - colorCode: the ANSI color code (e.g., "\033[31m" for red)
// - indexes: which character positions to color (e.g., [2, 3, 4])
//
// Returns: the colored ASCII art as a string
func RenderWithColor(input string, banner Banner, colorCode string, indexes []int) string {
	// Step 1: Decode any \n escape sequences to real newlines
	// Example: "Hello\nWorld" becomes actual two lines
	input = decodeEscapedNewlines(input)

	// Step 2: Handle empty input
	// If the input is empty, there's nothing to render
	if input == "" {
		return ""
	}

	// Step 3: Check if input contains newlines
	// If it does, we need to handle it differently (split by lines)
	if strings.Contains(input, "\n") {
		// Handle multi-line input with color
		return renderMultiLineWithColor(input, banner, colorCode, indexes)
	}

	// Step 4: Single line rendering (most common case)
	// This is where we do the character-by-character coloring
	return renderSingleLineWithColor(input, banner, colorCode, indexes)
}

// renderSingleLineWithColor renders a single line with colors
// This is the core coloring logic - character by character
func renderSingleLineWithColor(input string, banner Banner, colorCode string, indexes []int) string {
	// Constants
	const height = charHeight // 8 lines per character

	// Create a string builder for efficient string concatenation
	var builder strings.Builder

	// Create an empty glyph for unknown characters
	empty := make([]string, height)

	// Convert input string to slice of runes (characters)
	// This handles Unicode properly
	chars := []rune(input)

	// Step 1: Render row by row (each character is 8 rows tall)
	// We go through rows 0 to 7
	for row := 0; row < height; row++ {
		// Step 2: For each character in the input
		for charIndex, ch := range chars {
			// Step 3: Get the ASCII art for this character from the banner
			glyph, ok := banner[ch]
			if !ok {
				// Character not found in banner, use empty glyph
				glyph = empty
			}

			// Step 4: Check if this character needs to be colored
			// We check if this character's position is in our indexes slice
			shouldColor := ContainsIndex(indexes, charIndex)

			// Step 5: Add color code if needed (at the start of each row for this char)
			if shouldColor {
				// Add the color code BEFORE this character's row
				builder.WriteString(colorCode)
				// Add this character's current row
				builder.WriteString(glyph[row])
				// Add the reset code AFTER this character's row
				builder.WriteString(ResetColor)
			} else {
				// No color needed, just add the row normally
				builder.WriteString(glyph[row])
			}
		}

		// Step 6: End of this row, add newline
		builder.WriteRune('\n')
	}

	// Return the complete colored ASCII art
	return builder.String()
}

// renderMultiLineWithColor handles input with newlines
// Splits the input by newlines and renders each part
func renderMultiLineWithColor(input string, banner Banner, colorCode string, indexes []int) string {
	var builder strings.Builder

	// Split input by newlines
	lines := strings.Split(input, "\n")
	hadText := false

	// Keep track of character position across all lines
	// This is important for correct coloring when we have multiple lines
	totalPos := 0

	// Process each line
	for i, line := range lines {
		isLast := i == len(lines)-1

		if line != "" {
			// Non-empty line - render it with color

			// Calculate which indexes apply to THIS line
			// We need to adjust indexes based on totalPos
			lineIndexes := make([]int, 0)
			for _, idx := range indexes {
				// Check if this index falls within this line
				if idx >= totalPos && idx < totalPos+len(line) {
					// Adjust index to be relative to this line
					lineIndexes = append(lineIndexes, idx-totalPos)
				}
			}

			// Render this line with its adjusted indexes
			asciiBlock := renderSingleLineWithColor(line, banner, colorCode, lineIndexes)
			builder.WriteString(asciiBlock)
			hadText = true

			// Update total position (including the newline character)
			totalPos += len(line) + 1

		} else {
			// Empty line handling
			if !isLast {
				// Empty line in the middle, add a blank line
				builder.WriteRune('\n')
			} else if hadText {
				// Empty line at the end, keep the newline
				builder.WriteRune('\n')
			}
			// Update position for the newline
			totalPos++
		}
	}

	return builder.String()
}
