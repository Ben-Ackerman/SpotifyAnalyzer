package app

import (
	"regexp"
	"strings"
)

type Track struct {
	Artist    string
	Name      string
	GeniusURI string
	Lyrics    string
}

type lyricsCleaningOptions int

const (
	removeSectionHeaders lyricsCleaningOptions = iota
	trimWhiteSpace
)

const wordTrimChars = " \t\n,\"()[]#!?.-!~+^*/"

func cleanLyrics(lyrics string, opt ...lyricsCleaningOptions) string {
	for _, option := range opt {
		if option == removeSectionHeaders {
			lyrics = cleanLyricsMoveSectionHeader(lyrics)
		} else if option == trimWhiteSpace {
			lyrics = cleanLyricsTrimWhiteSpace(lyrics)
		}
	}
	return lyrics
}

func cleanLyricsTrimWhiteSpace(lyrics string) string {
	lines := strings.Split(lyrics, "\n")
	var sb strings.Builder
	for _, line := range lines {
		trimed := strings.TrimSpace(line)
		sb.WriteString(trimed)
		sb.WriteString("\n")
	}

	return sb.String()
}

func cleanLyricsMoveSectionHeader(lyrics string) string {
	lines := strings.Split(lyrics, "\n")
	re := regexp.MustCompile("[*]")
	var sb strings.Builder

	for _, line := range lines {
		sb.WriteString(re.ReplaceAllString(line, ""))
		sb.WriteString("\n")
	}

	return sb.String()
}

func getWordCounts(lyrics string) map[string]int {
	lines := strings.Split(lyrics, "\n")
	wordCount := make(map[string]int)
	for _, line := range lines {
		words := strings.Split(line, " ")
		for _, word := range words {
			word = strings.ToLower(word)
			word := strings.Trim(word, wordTrimChars)
			if len(word) == 0 {
				//Do not include the empty string
				continue
			}
			_, inMap := wordCount[word]
			if inMap {
				wordCount[word]++
			} else {
				wordCount[word] = 1
			}
		}
	}
	return wordCount
}

func getTopNWords(wordMap map[string]int, n int) []string {
	// Make a list of all words
	i := 0
	keys := make([]string, len(wordMap))
	for key := range wordMap {
		keys[i] = key
		i++
	}

	//Run bubble sort N times on keys to get the top N words
	//Sorting the words from largest to smallest
	for j := 0; j < n; j++ {
		for i := len(keys) - 1; i > 0; i-- {
			if wordMap[keys[i-1]] < wordMap[keys[i]] {
				keys[i], keys[i-1] = keys[i-1], keys[i]
			}
		}
	}

	return keys[0:n]
}
