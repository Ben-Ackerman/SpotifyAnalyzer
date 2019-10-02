package app

import (
	"regexp"
	"strings"
)

type lyricsCleaningOptions int

const (
	removeSectionHeaders lyricsCleaningOptions = iota
	trimWhiteSpace
)

const wordTrimChars = " \t\n,\"()[]#!?.-!~+^*/"

// cleanLyrics takes in a string representing lyrics and a list of options signifying what operations the caller wants performed on the lyrics.
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

// cleanLyricsTrimWhiteSpaces takes in a string representing lyrics and returns a new string where extra white space from each line has neen removed.
// This function also removes blank lines
func cleanLyricsTrimWhiteSpace(lyrics string) string {
	lines := strings.Split(lyrics, "\n")
	var sb strings.Builder
	for _, line := range lines {
		trimed := strings.TrimSpace(line)
		if len(trimed) > 0 {
			sb.WriteString(trimed)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// cleanLyricsMoveSectionHeader takes in a string representing lyrics and returns a new string where the lyric headers have been removed.
// An example of a lyrics header is [VERSE 1] or [CHORUS 2]
func cleanLyricsMoveSectionHeader(lyrics string) string {
	lines := strings.Split(lyrics, "\n")
	re := regexp.MustCompile(`\[.*\]`)
	var sb strings.Builder

	for _, line := range lines {
		sb.WriteString(re.ReplaceAllString(line, ""))
		sb.WriteString("\n")
	}

	return sb.String()
}

// getWordCounts takes in a string represneting lyrics and returns a map which maps a given word to the number of occurances of that word within the lyrics.
// When counting each word every word is changed to lower case and the characters " \t\n,\"()[]#!?.-!~+^*/" are trimed from each word.
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

// getTopNWords takes in a map of words to occurances and returns a list of the n most frequent words.
// The resulting list of words is sorted from most frequent to least frequent
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
	size := min(len(wordMap), n)
	for j := 0; j < size; j++ {
		for i := len(keys) - 1; i > 0; i-- {
			if wordMap[keys[i-1]] < wordMap[keys[i]] {
				keys[i], keys[i-1] = keys[i-1], keys[i]
			}
		}
	}

	return keys[0:size]
}

// Finds the min of two ints
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
