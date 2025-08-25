package utils

import (
	"os"
	"testing"
)

func TestLoadWords(t *testing.T) {
	// Create a temporary file with sample words
	tmpFile, err := os.CreateTemp("", "words_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after test

	content := "apple\nbanana\ncarrot\n\n  \ndate\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call the function
	words, err := LoadWords(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadWords failed: %v", err)
	}

	// Expected output
	expected := []string{"apple", "banana", "carrot", "date"}

	if len(words) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(words))
	}

	for i, word := range expected {
		if words[i] != word {
			t.Errorf("Expected word %q at index %d, got %q", word, i, words[i])
		}
	}
}
