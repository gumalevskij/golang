package main

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strconv"
	"xkcd-fetcher/pkg/database"
	"xkcd-fetcher/pkg/words"
	"xkcd-fetcher/pkg/xkcd"
)

type Config struct {
	SourceURL     string `yaml:"source_url"`
	DbFile        string `yaml:"db_file"`
	StopWordsFile string `yaml:"stopwords_file"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

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
		normComics[strconv.Itoa(comic.Num)] = map[string]interface{}{
			"url":      comic.Img,
			"keywords": result,
		}
	}

	return normComics, nil
}

func main() {
	outputFlag := flag.Bool("o", false, "Output to screen")
	limitFlag := flag.Int("n", -1, "Number of comics to fetch")
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
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
