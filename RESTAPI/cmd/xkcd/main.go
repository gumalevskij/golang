package main

import (
	"flag"
	"fmt"
	"log"
	"xkcd-fetcher/pkg/database"
	"xkcd-fetcher/pkg/words"
	"xkcd-fetcher/pkg/xkcd"
)

func NormalizeComicsNew(stopWordsFile string, comics []xkcd.Comic) (map[string]interface{}, error) {
	stopWords, err := words.ReadStopWords(stopWordsFile)
	if err != nil {
		fmt.Println("Error reading stop words file:", err)
		return nil, err
	}

	normComics := make(map[string]interface{})
	stemmer := words.NewStemmer()

	for _, comic := range comics {
		result := words.Normalize(stemmer, stopWords, comic.Transcript)
		normComics[fmt.Sprint(comic.Num)] = map[string]interface{}{
			"url":      comic.Img,
			"keywords": result,
		}
	}

	return normComics, nil
}

func main() {
	outputFlag := flag.Bool("o", false, "Output to screen")
	limitFlag := flag.Int("n", -1, "Number of comics to fetch")
	flag.Parse()

	config, err := xkcd.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	comics, err := xkcd.FetchComics(config.SourceURL, *limitFlag)
	if err != nil {
		log.Fatalf("Failed to fetch comics: %v", err)
	}

	normComics, err := NormalizeComicsNew(config.StopWordsFile, comics)
	if err != nil {
		log.Fatalf("Failed to normalize comics: %v", err)
	}

	if *outputFlag {
		fmt.Println(normComics)
	} else {
		if err := database.SaveComics(config.DbFile, normComics); err != nil {
			log.Fatalf("Failed to save comics: %v", err)
		}
	}
}
