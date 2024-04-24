package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"time"
	"xkcd-searcher/pkg/database"
	"xkcd-searcher/pkg/words"
	"xkcd-searcher/pkg/xkcd"
)

type ComicScore struct {
	ID    string
	Score int
	URL   string
}

func IndexingSearchComics(normalizedQuery []string, index database.SearchIndex, comicsDB database.Comics) []string {
	matchCounts := make(map[string]int)
	for _, word := range normalizedQuery {
		if ids, exists := index[word]; exists {
			for _, id := range ids {
				matchCounts[id]++
			}
		}
	}

	var comicsScores []ComicScore
	for id, count := range matchCounts {
		comicsScores = append(comicsScores, ComicScore{ID: id, Score: count, URL: comicsDB[id].Url})
	}

	sort.Slice(comicsScores, func(i, j int) bool {
		return comicsScores[i].Score > comicsScores[j].Score
	})

	var urls []string
	for i := 0; i < len(comicsScores) && i < 10; i++ {
		urls = append(urls, comicsScores[i].URL)
	}
	return urls
}

func InefficientSearchComics(normalizedQuery []string, comicsDB database.Comics) []string {
	queryWords := make(map[string]bool)
	for _, word := range normalizedQuery {
		queryWords[word] = true
	}

	matchCounts := make(map[string]int)
	for id, comic := range comicsDB {
		count := 0
		for _, word := range comic.Keywords {
			if queryWords[word] {
				count++
			}
		}
		if count > 0 {
			matchCounts[id] = count
		}
	}

	var comicsScores []ComicScore
	for id, count := range matchCounts {
		comicsScores = append(comicsScores, ComicScore{ID: id, Score: count, URL: comicsDB[id].Url})
	}

	sort.Slice(comicsScores, func(i, j int) bool {
		return comicsScores[i].Score > comicsScores[j].Score
	})

	var urls []string
	for i := 0; i < len(comicsScores) && i < 10; i++ {
		urls = append(urls, comicsScores[i].URL)
	}
	return urls
}

func containsAny(words []string, wordMap map[string]bool) bool {
	for _, word := range words {
		if wordMap[word] {
			return true
		}
	}
	return false
}

func main() {
	defer func(start time.Time) {
		fmt.Printf("%v\n", time.Since(start))
	}(time.Now())
	indexSearch := flag.Bool("i", false, "Switch on indexing search")
	inputQuery := flag.String("s", "", "Input string to query")
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	normComics, err := database.LoadComics(config.DbFile)
	if err != nil {
		normComics = make(database.Comics)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		if err := database.SaveComicsCache(config.DbFile, normComics); err != nil {
			log.Printf("Failed to save comics: %v", err)
		}
		stop()
	}()

	numJobs, err := xkcd.GetLastComicBinary(config.SourceURL, config.MinNumberComics, config.MaxNumberComics)
	if err != nil {
		log.Printf("Failed to find last comic: %v", err)
	}

	stopWords, err := words.ReadStopWords(config.StopWordsFile)
	if err != nil {
		log.Printf("Failed to read stop words file: %v", err)
	}

	stemmer := words.NewStemmer()
	jobs := make(chan int, numJobs)
	results := make(chan *xkcd.Comic, numJobs)
	for w := 1; w <= config.Parallel; w++ {
		go func(ctx context.Context) {
			for {
				select {
				case j, ok := <-jobs:
					if !ok {
						return
					}
					comic, err := xkcd.FetchComics(config.SourceURL, j)
					if err != nil {
						log.Printf("Failed to fetch comics %d: %v", j, err)
						results <- nil
						continue
					}
					if comic != nil {
						fullText := comic.Transcript + " " + comic.Alt + " " + comic.Title
						comic.NormalizedText = words.Normalize(stemmer, stopWords, fullText)
					}
					results <- comic
				case <-ctx.Done():
					stop()
					return
				}
			}
		}(ctx)
	}

	realCountJobs := 0
	for j := 1; j <= numJobs; j++ {
		if _, exists := normComics[strconv.Itoa(j)]; exists {
			continue
		}

		realCountJobs++
		jobs <- j
	}
	close(jobs)

	for doneJob := 1; doneJob <= realCountJobs; doneJob++ {
		comic := <-results
		if comic != nil {
			normComics[strconv.Itoa(comic.Num)] = database.NormalizedComic{
				Url:      comic.Img,
				Keywords: comic.NormalizedText,
			}
		}

		if err != nil {
			log.Fatalf("Failed to normalize comics: %v", err)
		}

		if doneJob%10 == 0 {
			if err := database.SaveComicsCache(config.DbFile, normComics); err != nil {
				log.Printf("Failed to periodically save comics: %v", err)
			}
		}
	}

	if len(*inputQuery) != 0 {
		NormalizedQuery := words.Normalize(stemmer, stopWords, *inputQuery)
		index := database.BuildIndex(normComics)

		var urls []string
		if *indexSearch {
			urls = IndexingSearchComics(NormalizedQuery, index, normComics)
		} else {
			urls = InefficientSearchComics(NormalizedQuery, normComics)
		}

		for _, url := range urls {
			fmt.Println(url)
		}
	}
}
