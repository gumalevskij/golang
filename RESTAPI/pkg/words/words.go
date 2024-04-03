package words

import (
	"bufio"
	"strings"
	"os"
	"unicode"
)

func ReadStopWords(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stopWords := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stopWords[line] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopWords, nil
}

func Normalize(stemmer Stemmer, stopWords map[string]bool, text string) []string {
	var result []string
	uniqueWords := make(map[string]bool)
	words := strings.Fields(text)
	for _, word := range words {
		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		
		stemmedWord := stemmer.Stem(word)
		stemmedWord = strings.Trim(stemmedWord, "’")
		if _, exists := uniqueWords[stemmedWord]; !stopWords[word] && !exists {
			uniqueWords[stemmedWord] = true
			result = append(result, stemmedWord)
		}
	}
	return result
}
