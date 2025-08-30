package game

import (
	"path"
	"testing"

	"github.com/tomlaws/wordle/pkg/utils"
)

func TestWordListInitialization(t *testing.T) {
	wordList, err := NewWordList(path.Join(utils.Root, "assets", "words.txt"))
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

func TestRandomWord(t *testing.T) {
	wordList, err := NewWordList(path.Join(utils.Root, "assets", "words.txt"))
	if err != nil {
		t.Fatalf("Failed to load word list: %v", err)
	}
	word := wordList.RandomWord()
	if word == "" {
		t.Errorf("RandomWord returned empty string")
	}
	found := false
	for _, w := range wordList.words {
		if w == word {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("RandomWord returned a word not in the list: %s", word)
	}
}

func TestIsValidWord(t *testing.T) {
	wordList, err := NewWordList(path.Join(utils.Root, "assets", "words.txt"))
	if err != nil {
		t.Fatalf("Failed to load word list: %v", err)
	}
	validWords := []string{"apple", "magic", "TaBlE"}
	invalidWords := []string{"xyzzy", "foobar", "qwerty"}

	for _, word := range validWords {
		if !wordList.IsValidWord(word) {
			t.Errorf("Expected %s to be a valid word", word)
		}
	}

	for _, word := range invalidWords {
		if wordList.IsValidWord(word) {
			t.Errorf("Expected %s to be an invalid word", word)
		}
	}
}
