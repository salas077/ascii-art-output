package main

import (
	"fmt"
	"os"
)

func main() {
	// Step 1: Parse the command line arguments
	// This function reads os.Args and figures out:
	// - Do we need color? (UseColor)
	// - What color? (Color)
	// - What substring to color? (Substring)
	// - What text to render? (Text)
	// - What banner to use? (Banner)
	opts, err := ParseColorArgs(os.Args)

	// Step 2: Check if parsing failed
	// If there was an error (wrong format, missing args, etc.)
	// we show the usage message and exit
	if err != nil {
		// This is the EXACT usage message required by the audit
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println()
		fmt.Println("EX: go run . --output=<fileName.txt> something standard")
		return // Exit the program
	}

	// Step 3: Determine which banner file to load
	// Based on the banner name (standard, shadow, or thinkertoy)
	// we set the correct file path
	var bannerPath string
	switch opts.Banner {
	case "standard":
		// User wants standard banner
		bannerPath = "banners/standard.txt"
	case "shadow":
		// User wants shadow banner
		bannerPath = "banners/shadow.txt"
	case "thinkertoy":
		// User wants thinkertoy banner
		bannerPath = "banners/thinkertoy.txt"
	default:
		// User provided an invalid banner name
		// Show error message with what they typed and what's available
		fmt.Printf("Error: Invalid banner '%s'\n", opts.Banner)
		fmt.Println("Available banners: standard, shadow, thinkertoy")
		return // Exit the program
	}

	// Step 4: Load the banner file from disk
	// This reads the file and creates the Banner map
	// where each character maps to its ASCII art
	banner, err := LoadBanner(bannerPath)

	// Step 5: Check if loading the banner failed
	// Could fail if file doesn't exist, is corrupted, etc.
	if err != nil {
		// Show which file failed and why
		fmt.Printf("Error: Could not load banner file '%s'\n", bannerPath)
		fmt.Printf("Details: %v\n", err)
		return // Exit the program
	}

	// Step 6: Decide whether to render with color or without color
	// This is the main branching point in our program
	if opts.UseColor {
		// ===== COLOR MODE =====
		// User wants colored output

		// Step 6a: Validate the color name
		// Check if the color exists in our colorCodes map
		// GetColorCode returns the ANSI code and a boolean (exists or not)
		colorCode, exists := GetColorCode(opts.Color)

		// Step 6b: Check if the color is valid
		if !exists {
			// User typed an invalid color name (e.g., "pink")
			// Show error with what they typed and available colors
			fmt.Printf("Error: Invalid color '%s'\n", opts.Color)
			fmt.Println("Available colors: red, green, yellow, blue, magenta, cyan, white, orange, black")
			return // Exit the program
		}

		// Step 6c: Check if substring is empty
		// If user provided an empty substring, we show a warning and render normally
		if opts.Substring == "" && len(os.Args) > 3 {
			// User explicitly gave an empty substring like: --color=red "" "text"
			// This is probably a mistake, so we warn them
			fmt.Println("Warning: Empty substring provided. Rendering without color.")
			fmt.Println()
			// Render normally without color
			output := RenderInput(opts.Text, banner)
			if opts.OutputFile != "" {
				// Write to file
				err := os.WriteFile(opts.OutputFile, []byte(output), 0644)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			} else {
				// Write to screen
				fmt.Print(output)
			}
			return
		}

		// Step 6d: Find which character positions need to be colored
		// This returns a slice of integers representing positions
		// Example: "a kitten" with substring "kit" returns [2, 3, 4]
		// If substring is empty (no substring arg), it returns ALL positions (color everything)
		indexes := FindSubstringIndexes(opts.Text, opts.Substring)

		// Step 6e: Render the text with colors
		// This is our NEW function that we created
		// It renders character-by-character and adds color codes where needed
		output := RenderWithColor(opts.Text, banner, colorCode, indexes)

		// Step 6f: Print the colored output
		if opts.OutputFile != "" {
			// Write to file
			err := os.WriteFile(opts.OutputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else {
			// Write to screen
			fmt.Print(output)
		}

	} else {
		// ===== NORMAL MODE (NO COLOR) =====
		// User didn't specify --color flag
		// Render normally, just like the old ascii-art program

		// Step 6g: Render the text without colors
		// This is our existing function from before
		output := RenderInput(opts.Text, banner)

		// Step 6h: Print the normal output
		if opts.OutputFile != "" {
			// Write to file
			err := os.WriteFile(opts.OutputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else {
			// Write to screen
			fmt.Print(output)
		}
	}

	// Program ends successfully
	// No need for explicit return at the end of main
}
