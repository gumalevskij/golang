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

func SaveComics(path string, comics Comics) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
