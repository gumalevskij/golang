package database

import (
	"encoding/json"
	"io/ioutil"
	"gopkg.in/yaml.v2"
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

func SaveComics(path string, comics map[string]interface{}) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
