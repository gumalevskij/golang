package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"xkcd-fetcher/pkg/database"
	"xkcd-fetcher/pkg/words"
	"xkcd-fetcher/pkg/xkcd"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL     string `yaml:"source_url"`
	DbFile        string `yaml:"db_file"`
	StopWordsFile string `yaml:"stopwords_file"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

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

func getMaxComicNum(comics database.Comics) (int, error) {
	maxNum := 0
	for numStr := range comics {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, err
		}
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum, nil
}

func main() {
	outputFlag := flag.Bool("o", false, "Output to screen")
	limitFlag := flag.Int("n", -1, "Number of comics to fetch, -1 for all available")
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var normComics database.Comics
	var comics []xkcd.Comic

	normComics, err = database.LoadComics(config.DbFile)
	if err != nil || len(normComics) == 0 {
		comics, err = xkcd.FetchComics(config.SourceURL, 0, *limitFlag)
		if err != nil {
			log.Fatalf("Failed to fetch comics: %v", err)
		}
		normComics = make(database.Comics)
	} else {
		maxComicsNum, err := getMaxComicNum(normComics)
		if err != nil {
			log.Fatalf("Failed to find max comic number: %v", err)
			return
		}

		if *limitFlag != -1 {
			*limitFlag = maxComicsNum + *limitFlag
		}
		comics, err = xkcd.FetchComics(config.SourceURL, maxComicsNum, *limitFlag)
		if err != nil {
			log.Fatalf("Failed to fetch comics: %v", err)
		}
	}

	normComics, err = NormalizeComics(config.StopWordsFile, comics, normComics)
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
