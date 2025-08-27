package game

import (
	"github.com/tomlaws/wordle/pkg/utils"
)

func NewWordList(path string) (*WordList, error) {
	words, err := utils.LoadWords(path)
	if err != nil {
		return nil, err
	}
	index := make(map[string]int)
	for i, word := range words {
		index[word] = i
	}
	return &WordList{words: words, index: index}, nil
}

func (wl *WordList) RandomWord() string {
	if len(wl.words) == 0 {
		return ""
	}
	return wl.words[utils.RandomInt(0, len(wl.words)-1)]
}

func (wl *WordList) IsValidWord(word string) bool {
	_, exists := wl.index[word]
	return exists
}
