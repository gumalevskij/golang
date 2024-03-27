package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "unicode"
)

func readStopWords(filename string) (map[string]bool, error) {
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

func main() {
    var input string
    flag.StringVar(&input, "s", "", "Input string to normalize")
    flag.Parse()

    stopWords, err := readStopWords("stopWordsEnglish.txt")
    if err != nil {
        fmt.Println("Error reading stop words file:", err)
        return
    }

    words := strings.Fields(input)
    uniqueWords := make(map[string]bool)
    var result []string

	stemmer := NewStemmer()

    for _, word := range words {
		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
            return !unicode.IsLetter(r) && !unicode.IsNumber(r)
        })
		stemmedWord := stemmer.Stem(word);
		if _, exists := uniqueWords[stemmedWord]; !stopWords[word] && !exists {
			uniqueWords[stemmedWord] = true
			result = append(result, stemmedWord)
		}
	}

    fmt.Println(strings.Join(result, " "))
}
