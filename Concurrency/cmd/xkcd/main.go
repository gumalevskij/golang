package main

import (
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
	start := time.Now()
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		if err := database.SaveComicsCaсhe(config.DbFile, normComics); err != nil {
			log.Printf("Failed to save comics: %v", err)
		}
		os.Exit(0)
	}()

	low, high := 2917, 2917 // or 0 to 5000
	numJobs, err := xkcd.GetLastComicBinary(config.SourceURL, low, high)

	stopWords, err := words.ReadStopWords(config.StopWordsFile)
	if err != nil {
		log.Printf("Failed to read stop words file: %v", err)
	}

	stemmer := words.NewStemmer()
	jobs := make(chan int, numJobs)
	results := make(chan []xkcd.Comic, numJobs)
	for w := 1; w <= config.Parallel; w++ {
		go func() {
			for j := range jobs {
				comics, err := xkcd.FetchComics(config.SourceURL, j-1, j)
				for i := range comics {
					fullText := comics[i].Transcript + " " + comics[i].Alt
					comics[i].NormalizedText = words.Normalize(stemmer, stopWords, fullText)
				}
				if err != nil {
					log.Fatalf("Failed to fetch comics: %v", err)
				}
				results <- comics
			}
		}()
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
		comics := <-results
		for _, comic := range comics {
			normComics[strconv.Itoa(comic.Num)] = database.NormalizedComic{
				Url:      comic.Img,
				Keywords: comic.NormalizedText,
			}
		}
		if err != nil {
			log.Fatalf("Failed to normalize comics: %v", err)
		}

		if doneJob%10 == 0 {
			if err := database.SaveComicsCaсhe(config.DbFile, normComics); err != nil {
				log.Printf("Failed to periodically save comics: %v", err)
			}
		}
	}

	if err := database.SaveComicsCaсhe(config.DbFile, normComics); err != nil {
		log.Fatalf("Failed to save comics: %v", err)
	}

	duration := time.Since(start)

	fmt.Println(duration.Seconds())
}
