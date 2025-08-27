package game

import (
	"testing"
)

func TestGameInitialization(t *testing.T) {
	maxAttempts := 6
	game := NewGame("apple", maxAttempts)
	if game.Answer != "apple" {
		t.Errorf("Expected answer to be 'apple', got %s", game.Answer)
	}
	if game.MaxAttempts != maxAttempts {
		t.Errorf("Expected max attempts to be %d, got %d", maxAttempts, game.MaxAttempts)
	}
	// Check size of Attempts slice capacity
	if cap(game.Attempts) != maxAttempts {
		t.Errorf("Expected attempts capacity to be %d, got %d", maxAttempts, cap(game.Attempts))
	}
	// Check initial state
	if len(game.Attempts) != 0 {
		t.Errorf("Expected 0 attempts at start, got %d", len(game.Attempts))
	}
	if game.State != InProgress {
		t.Errorf("Expected initial game state to be InProgress, got %v", game.State)
	}
}

func TestMakeGuess_CorrectGuess(t *testing.T) {
	game := NewGame("apple", 6)
	result, err := game.MakeGuess("apple")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if game.State != Won {
		t.Errorf("Expected game state to be Won, got %v", game.State)
	}
	for i, lr := range result {
		if lr.MatchType != Hit {
			t.Errorf("Expected Hit at position %d, got %v", i, lr.MatchType)
		}
	}
}

func TestMakeGuess_IncorrectGuess(t *testing.T) {
	game := NewGame("apple", 6)
	result, err := game.MakeGuess("grape")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if game.State != InProgress {
		t.Errorf("Expected game state to be InProgress, got %v", game.State)
	}
	// result should not be all Hits
	allHits := true
	for _, lr := range result {
		if lr.MatchType != Hit {
			allHits = false
			break
		}
	}
	if allHits {
		t.Errorf("Expected not all Hits for incorrect guess, got all Hits")
	}
}

func TestMakeGuess_PartialMatch(t *testing.T) {
	game := NewGame("apple", 6)
	result, err := game.MakeGuess("grape")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if game.State != InProgress {
		t.Errorf("Expected game state to be InProgress, got %v", game.State)
	}
	// result should be Miss, Miss, Present, Present, Hit
	// "grape" vs "apple"
	expected := []MatchType{Miss, Miss, Present, Present, Hit}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_InvalidLength(t *testing.T) {
	game := NewGame("apple", 6)
	result, err := game.MakeGuess("app")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for invalid guess length, got %v", result)
	}
}

func TestMakeGuess_MaxAttempts(t *testing.T) {
	game := NewGame("apple", 2)
	game.MakeGuess("grape")
	game.MakeGuess("grape")
	if game.State != Lost {
		t.Errorf("Expected game state to be Lost after max attempts, got %v", game.State)
	}
}

func TestMakeGuess_GameOver(t *testing.T) {
	game := NewGame("apple", 1)
	game.MakeGuess("grape")
	if game.State != Lost {
		t.Errorf("Expected game state to be Lost after first guess, got %v", game.State)
	}
	result, err := game.MakeGuess("apple")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result after game over, got %v", result)
	}
}

func TestMakeGuess_HitTakesPrecedenceOverPresent(t *testing.T) {
	game := NewGame("kitty", 6)
	result, _ := game.MakeGuess("empty")
	// "kitty" vs "empty"
	// k: Miss, e: Miss, t: Miss, t: Hit, y: Hit
	expected := []MatchType{Miss, Miss, Miss, Hit, Hit}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_ExtraGuessLetterAfterHitIsMiss(t *testing.T) {
	game := NewGame("smile", 6)
	result, _ := game.MakeGuess("skill")
	// "skill" vs "smile"
	// s: Hit, k: Miss, i: Hit, l: Hit, l: Miss (the last one should be Miss instead of Present)
	expected := []MatchType{Hit, Miss, Hit, Hit, Miss}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_ExtraGuessLetterAfterPresentIsMiss(t *testing.T) {
	game := NewGame("smile", 6)
	result, _ := game.MakeGuess("alley")
	// "alley" vs "smile"
	// a: Miss, l: Present, l: Miss, e: Present, y: Miss (the second l should be Miss instead of Present)
	expected := []MatchType{Miss, Present, Miss, Present, Miss}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_CaseInsensitivity(t *testing.T) {
	game := NewGame("Apple", 6)
	result, _ := game.MakeGuess("aPpLe")
	// "aPpLe" vs "Apple"
	// a: Hit, P: Hit, p: Hit, L: Hit, e: Hit
	expected := []MatchType{Hit, Hit, Hit, Hit, Hit}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_CaseInsensitivityWithPresents(t *testing.T) {
	game := NewGame("Apple", 6)
	result, _ := game.MakeGuess("pPale")
	// "pPale" vs "Apple"
	// p: Present, P: Hit, a: Present, l: Hit, e: Hit
	expected := []MatchType{Present, Hit, Present, Hit, Hit}
	for i, lr := range result {
		if lr.MatchType != expected[i] {
			t.Errorf("At pos %d: expected %v, got %v", i, expected[i], lr.MatchType)
		}
	}
}

func TestMakeGuess_MakeGuessAfterWin(t *testing.T) {
	game := NewGame("apple", 6)
	game.MakeGuess("apple")
	if game.State != Won {
		t.Errorf("Expected game state to be Won, got %v", game.State)
	}
	result, err := game.MakeGuess("grape")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result after game won, got %v", result)
	}
}
