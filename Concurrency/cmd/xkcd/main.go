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

func NormalizeComics(stopWordsFile string, comics []xkcd.Comic, normComics database.Comics) (database.Comics, error) {
	stopWords, err := words.ReadStopWords(stopWordsFile)
	if err != nil {
		fmt.Println("Error reading stop words file:", err)
		return normComics, err
	}

	stemmer := words.NewStemmer()

	for _, comic := range comics {
		fullText := comic.Transcript + " " + comic.Alt
		result := words.Normalize(stemmer, stopWords, fullText)
		(normComics)[strconv.Itoa(comic.Num)] = database.NormalizedComic{
			Url:      comic.Img,
			Keywords: result,
		}
	}

	return normComics, nil
}

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
		if err := database.SaveComics(config.DbFile, normComics); err != nil {
			log.Printf("Failed to save comics: %v", err)
		}
		os.Exit(0)
	}()

	low, high := 2917, 2917 // or 0 to 5000
	numJobs, err := xkcd.GetLastComicBinary(config.SourceURL, low, high)

	jobs := make(chan int, numJobs)
	results := make(chan []xkcd.Comic, numJobs)
	for w := 1; w <= config.Parallel; w++ {
		go func() {
			for j := range jobs {
				comics, err := xkcd.FetchComics(config.SourceURL, j-1, j)
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

	var a int
	for a = 1; a <= realCountJobs; a++ {
		comics := <-results
		normComics, err = NormalizeComics(config.StopWordsFile, comics, normComics)
		if err != nil {
			log.Fatalf("Failed to normalize comics: %v", err)
		}

		if a%10 == 0 {
			if err := database.SaveComics(config.DbFile, normComics); err != nil {
				log.Printf("Failed to periodically save comics: %v", err)
			}
		}
	}

	if err := database.SaveComics(config.DbFile, normComics); err != nil {
		log.Fatalf("Failed to save comics: %v", err)
	}

	duration := time.Since(start)

	fmt.Println(duration.Seconds())
}
