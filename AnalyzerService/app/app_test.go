package app

import (
	"testing"
)

func TestCleanLyrics1(t *testing.T) {
	testLyrics := "[TITLE]\n[VERSE] \n [VERSE]"
	result := cleanLyrics(testLyrics, removeSectionHeaders, trimWhiteSpace)
	if len(result) != 0 {
		t.Errorf("FAILED: result was non-zero\nresult = \"%s\"", result)
	}
}

func TestCleanLyrics2(t *testing.T) {
	testLyrics := "[TITLE]\nline1\nline2 \n [VERSE]"
	result := cleanLyrics(testLyrics, removeSectionHeaders, trimWhiteSpace)
	if result != "line1\nline2\n" {
		t.Errorf("FAILED: result was non-zero\nresult = \"%s\"", result)
	}
}

func TestGetWordCount1(t *testing.T) {
	testLyrics := "Hello?   \n hello  \tworld, ,"
	wordCounts := getWordCounts(testLyrics)
	if len(wordCounts) != 2 {
		t.Errorf("FAILED: incorrect number of words, wordCounts = %v", wordCounts)
	}
	if wordCounts["hello"] != 2 {
		t.Errorf("FAILED: incorrect number of hello found. Should be 2 found %d", wordCounts["hello"])
	}
	if wordCounts["world"] != 1 {
		t.Errorf("FAILED: incorrect number of world found. Should be 1 found %d", wordCounts["hello"])
	}
}

func TestGetTopNWords(t *testing.T) {
	wordCounts := map[string]int{
		"low":    3,
		"yes":    50,
		"no":     51,
		"maybe":  48,
		"should": 3,
		"here":   75,
	}
	result := getTopNWords(wordCounts, 2)
	if len(result) != 2 {
		t.Errorf("FAILED: To many words returned should have returned 2 words instead returned %d words", len(result))
	}
	if result[0] != "here" || result[1] != "no" {
		t.Errorf("FAILED: incorrect words returned should be [here, no] insead is %v", result)
	}
}

func TestStopWords(t *testing.T) {
	s := &Server{}
	s.InitStopWords()
	wordCount := map[string]int{
		"a":       1,
		"an":      3,
		"asdfasd": 5,
	}
	result := s.removeStopWords(wordCount)
	if len(result) != 1 {
		t.Errorf("FAILED, error removing stop words final count is not correct.  len = %d", len(result))
	}
}
