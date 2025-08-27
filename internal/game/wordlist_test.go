package game

import (
	"path"
	"testing"
)

func TestWordListInitialization(t *testing.T) {
	wordList, err := NewWordList(path.Join("..", "..", "assets", "words.txt"))
	if err != nil {
		t.Fatalf("Failed to load word list: %v", err)
	}
	expectedWords := []string{"apple", "magic", "table", "candy", "ghost"}

	// Verify that each expected word exists in the index and matches the word at its index
	for _, word := range expectedWords {
		index, exists := wordList.index[word]
		if !exists {
			t.Errorf("Word %s not found in index", word)
			continue
		}
		if wordList.words[index] != word {
			t.Errorf("Expected word %s at index %d, got %s", word, index, wordList.words[index])
		}
	}
}

func TestWordListInitialization_NonExistentFile(t *testing.T) {
	_, err := NewWordList("non_existent_file.txt")
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}
