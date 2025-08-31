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
	Letter    rune      `json:"letter"`
	Position  int       `json:"position"`
	MatchType MatchType `json:"match_type"`
}

type Game struct {
	Answer     string
	MaxGuesses int
	Attempts   [][]LetterResult
	State      GameState
}

type WordList struct {
	words []string
	index map[string]int
}
