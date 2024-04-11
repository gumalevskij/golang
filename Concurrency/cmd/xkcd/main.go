package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
	"xkcd-fetcher/pkg/database"
	"xkcd-fetcher/pkg/words"
	"xkcd-fetcher/pkg/xkcd"
)

func main() {
	defer func(start time.Time) {
		fmt.Printf("%v\n", time.Since(start))
	}(time.Now())
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var normComics database.Comics
	normComics, err = database.LoadComics(config.DbFile)
	if err != nil || len(normComics) == 0 {
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
						log.Printf("Failed to fetch comics: %v", err)
						continue
					}
					if comic != nil {
						fullText := comic.Transcript + " " + comic.Alt
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

	if err := database.SaveComicsCache(config.DbFile, normComics); err != nil {
		log.Fatalf("Failed to save comics: %v", err)
	}
}
