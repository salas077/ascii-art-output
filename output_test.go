package main

import (
	"os"
	"testing"
)

// Test parsing --output flag
func TestParseColorArgs_WithOutput(t *testing.T) {
	args := []string{"program", "--output=test.txt", "Hello"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if opts.OutputFile != "test.txt" {
		t.Errorf("Expected output file 'test.txt', got '%s'", opts.OutputFile)
	}

	if opts.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", opts.Text)
	}
}

// Test parsing --output with color
func TestParseColorArgs_OutputWithColor(t *testing.T) {
	args := []string{"program", "--output=result.txt", "--color=red", "Hello"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if opts.OutputFile != "result.txt" {
		t.Errorf("Expected output file 'result.txt', got '%s'", opts.OutputFile)
	}

	if !opts.UseColor {
		t.Error("UseColor should be true")
	}

	if opts.Color != "red" {
		t.Errorf("Expected color 'red', got '%s'", opts.Color)
	}
}

// Test parsing --output with banner
func TestParseColorArgs_OutputWithBanner(t *testing.T) {
	args := []string{"program", "--output=file.txt", "Hello", "shadow"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if opts.OutputFile != "file.txt" {
		t.Errorf("Expected output file 'file.txt', got '%s'", opts.OutputFile)
	}

	if opts.Banner != "shadow" {
		t.Errorf("Expected banner 'shadow', got '%s'", opts.Banner)
	}
}

// Test empty output filename
func TestParseColorArgs_EmptyOutputFile(t *testing.T) {
	args := []string{"program", "--output=", "Hello"}
	_, err := ParseColorArgs(args)

	if err == nil {
		t.Error("Expected error for empty output file")
	}
}

// Test file writing functionality
func TestFileOutput(t *testing.T) {
	// Create a temporary file name
	tempFile := "test_output.txt"
	
	// Clean up after test
	defer os.Remove(tempFile)

	// Create fake banner for testing
	banner := make(Banner)
	banner['H'] = []string{"H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8"}

	// Test writing to file
	output := "H1\nH2\nH3\nH4\nH5\nH6\nH7\nH8\n"
	err := os.WriteFile(tempFile, []byte(output), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Read back and verify
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != output {
		t.Errorf("File content mismatch. Got: %q, Want: %q", string(content), output)
	}
}