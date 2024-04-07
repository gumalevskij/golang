package database

import (
	"encoding/json"
	"os"
)

type NormalizedComic struct {
    Url       string   `json:"url"`
    Keywords  []string `json:"keywords"`
}

type Comics map[string]NormalizedComic

func LoadComics(path string) (Comics, error) {
	var comics Comics

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return comics, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &comics)
	if err != nil {
		return nil, err
	}

	return comics, nil
}

func SaveComics(path string, comics Comics) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
