package game

type GameState int

const (
	InProgress GameState = iota
	Won
	Lost
)

type MatchType int

const (
	Miss MatchType = iota
	Present
	Hit
)

type LetterResult struct {
	Letter    rune
	Position  int
	MatchType MatchType
}

type Game struct {
	Answer      string
	MaxAttempts int
	Attempts    [][]LetterResult
	State       GameState
}

type WordList struct {
	words []string
	index map[string]int
}
