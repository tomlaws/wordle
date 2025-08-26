package game

import (
	"testing"
)

func TestWordListInitialization(t *testing.T) {
	wordList, err := NewWordList("../../assets/words.txt")
	if err != nil {
		t.Fatalf("Failed to load word list: %v", err)
	}
	expectedWords := []string{"apple", "magic", "table", "candy", "ghost"}

	// get the index first, then check the words at those indices
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
