package main

import "testing"

// Test that valid color names return correct ANSI codes
func TestGetColorCode_ValidColors(t *testing.T) {
	// Test a few basic colors
	code, ok := GetColorCode("red")
	if !ok || code != "\033[31m" {
		t.Errorf("GetColorCode(\"red\") failed")
	}

	code, ok = GetColorCode("blue")
	if !ok || code != "\033[34m" {
		t.Errorf("GetColorCode(\"blue\") failed")
	}
}

// Test that invalid color names are rejected
func TestGetColorCode_InvalidColor(t *testing.T) {
	_, ok := GetColorCode("pink")
	if ok {
		t.Error("GetColorCode should reject invalid color 'pink'")
	}
}

// Test that color names are case-insensitive
func TestGetColorCode_CaseInsensitive(t *testing.T) {
	code1, _ := GetColorCode("RED")
	code2, _ := GetColorCode("red")
	if code1 != code2 {
		t.Error("GetColorCode should be case-insensitive")
	}
}

// Test finding a substring that appears once
func TestFindSubstringIndexes_SingleMatch(t *testing.T) {
	indexes := FindSubstringIndexes("a kitten", "kit")
	expected := []int{2, 3, 4}

	if !equalSlices(indexes, expected) {
		t.Errorf("Expected %v, got %v", expected, indexes)
	}
}

// Test finding a substring that appears multiple times
func TestFindSubstringIndexes_MultipleMatches(t *testing.T) {
	indexes := FindSubstringIndexes("banana", "a")
	expected := []int{1, 3, 5}

	if !equalSlices(indexes, expected) {
		t.Errorf("Expected %v, got %v", expected, indexes)
	}
}

// Test when substring doesn't exist
func TestFindSubstringIndexes_NoMatch(t *testing.T) {
	indexes := FindSubstringIndexes("hello", "xyz")

	if len(indexes) != 0 {
		t.Errorf("Expected empty slice, got %v", indexes)
	}
}

// Test that empty substring returns all positions
func TestFindSubstringIndexes_EmptySubstring(t *testing.T) {
	indexes := FindSubstringIndexes("Hi", "")
	expected := []int{0, 1}

	if !equalSlices(indexes, expected) {
		t.Errorf("Expected %v, got %v", expected, indexes)
	}
}

// Test basic argument parsing without color
func TestParseColorArgs_NoColor(t *testing.T) {
	args := []string{"program", "Hello"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if opts.UseColor {
		t.Error("UseColor should be false")
	}

	if opts.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", opts.Text)
	}
}

// Test argument parsing with color flag
func TestParseColorArgs_WithColor(t *testing.T) {
	args := []string{"program", "--color=red", "Hello"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !opts.UseColor {
		t.Error("UseColor should be true")
	}

	if opts.Color != "red" {
		t.Errorf("Expected color 'red', got '%s'", opts.Color)
	}

	if opts.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", opts.Text)
	}
}

// Test argument parsing with substring
func TestParseColorArgs_WithSubstring(t *testing.T) {
	args := []string{"program", "--color=blue", "kit", "kitten"}
	opts, err := ParseColorArgs(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if opts.Substring != "kit" {
		t.Errorf("Expected substring 'kit', got '%s'", opts.Substring)
	}

	if opts.Text != "kitten" {
		t.Errorf("Expected text 'kitten', got '%s'", opts.Text)
	}
}

// Helper function to compare int slices
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
