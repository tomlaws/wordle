package game

import (
	"errors"
	"unicode"
)

func NewGame(answer string, maxAttempts int) *Game {
	return &Game{
		Answer:      answer,
		MaxAttempts: maxAttempts,
		Attempts:    make([][]LetterResult, 0, maxAttempts),
		State:       InProgress,
	}
}

func (g *Game) MakeGuess(guess string) ([]LetterResult, error) {
	if g.State != InProgress {
		return nil, errors.New("game is not in progress")
	}
	if len(guess) != len(g.Answer) {
		return nil, errors.New("invalid guess length")
	}
	result := make([]LetterResult, len(guess))
	answerRunes := []rune(g.Answer)
	guessRunes := []rune(guess)
	used := make([]bool, len(g.Answer))
	// First pass: check for hits
	for i, r := range guessRunes {
		if unicode.ToLower(r) == unicode.ToLower(answerRunes[i]) {
			result[i] = LetterResult{Letter: r, Position: i, MatchType: Hit}
			used[i] = true
		}
	}
	// Second pass: check for presents and misses
	for i, r := range guessRunes {
		if result[i].MatchType == Hit {
			continue
		}
		found := false
		for j, ar := range answerRunes {
			if !used[j] && unicode.ToLower(r) == unicode.ToLower(ar) {
				found = true
				used[j] = true
				break
			}
		}
		if found {
			result[i] = LetterResult{Letter: r, Position: i, MatchType: Present}
		} else {
			result[i] = LetterResult{Letter: r, Position: i, MatchType: Miss}
		}
	}
	g.Attempts = append(g.Attempts, result)
	if guess == g.Answer {
		g.State = Won
	} else if len(g.Attempts) >= g.MaxAttempts {
		g.State = Lost
	} else {
		g.State = InProgress
	}
	return result, nil
}
