package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL       string `yaml:"source_url"`
	DbFile          string `yaml:"db_file"`
	StopWordsFile   string `yaml:"stopwords_file"`
	Parallel        int    `yaml:"parallel"`
	MinNumberComics int    `yaml:"min_number_comics"`
	MaxNumberComics int    `yaml:"max_number_comics"`
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
