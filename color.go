package main

import (
	"fmt"
	"strings"
)

// ANSI color codes map
var colorCodes = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"orange":  "\033[38;5;208m",
	"black":   "\033[30m",
}

// Reset code
const ResetColor = "\033[0m"

// GetColorCode returns the ANSI code for a color name
func GetColorCode(colorName string) (string, bool) {
	colorName = strings.ToLower(colorName)
	code, exists := colorCodes[colorName]
	return code, exists
}

// FindSubstringIndexes finds all character positions of substring in text
func FindSubstringIndexes(text, substring string) []int {
	if substring == "" {
		// Color everything
		result := make([]int, len(text))
		for i := range text {
			result[i] = i
		}
		return result
	}

	var indexes []int
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			for j := 0; j < len(substring); j++ {
				indexes = append(indexes, i+j)
			}
		}
	}
	return indexes
}

// ContainsIndex checks if a number is in a slice
func ContainsIndex(indexes []int, target int) bool {
	for _, idx := range indexes {
		if idx == target {
			return true
		}
	}
	return false
}

// ColorOptions holds parsed arguments
type ColorOptions struct {
	UseColor   bool
	Color      string
	Substring  string
	Text       string
	Banner     string
	OutputFile string
}

// ParseColorArgs parses command line arguments
// Simple version - easier to understand!
func ParseColorArgs(args []string) (ColorOptions, error) {
	opts := ColorOptions{
		UseColor: false,
		Banner:   "standard",
	}

	// args[0] is program name
	if len(args) < 2 {
		return opts, fmt.Errorf("missing arguments")
	}

	// Start from index 1
	i := 1

	// Check for flags (--color= or --output=)
	for i < len(args) && strings.HasPrefix(args[i], "--") {
		if strings.HasPrefix(args[i], "--color=") {
			opts.UseColor = true
			opts.Color = args[i][8:] // After "--color="
			if opts.Color == "" {
				return opts, fmt.Errorf("empty color")
			}
		} else if strings.HasPrefix(args[i], "--output=") {
			opts.OutputFile = args[i][9:] // After "--output="
			if opts.OutputFile == "" {
				return opts, fmt.Errorf("empty output file")
			}
		} else {
			return opts, fmt.Errorf("unknown flag")
		}
		i++
	}

	// Need at least text after flags
	if i >= len(args) {
		return opts, fmt.Errorf("missing text")
	}

	// Now we have remaining arguments
	remaining := args[i:]

	switch len(remaining) {
	case 1:
		// Just text
		opts.Text = remaining[0]

	case 2:
		// Could be: [text, banner] OR [substring, text]
		if remaining[1] == "standard" || remaining[1] == "shadow" || remaining[1] == "thinkertoy" {
			// [text, banner]
			opts.Text = remaining[0]
			opts.Banner = remaining[1]
		} else if opts.UseColor {
			// [substring, text] - only valid with color flag
			opts.Substring = remaining[0]
			opts.Text = remaining[1]
		} else {
			return opts, fmt.Errorf("invalid arguments")
		}

	case 3:
		// [substring, text, banner] - only valid with color flag
		if !opts.UseColor {
			return opts, fmt.Errorf("too many arguments")
		}
		opts.Substring = remaining[0]
		opts.Text = remaining[1]
		opts.Banner = remaining[2]

	default:
		return opts, fmt.Errorf("invalid number of arguments")
	}

	return opts, nil
}
