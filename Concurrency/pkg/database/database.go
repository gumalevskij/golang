package database

import (
	"encoding/json"
	"os"
)

type NormalizedComic struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type Comics map[string]NormalizedComic

func LoadComics(path string) (Comics, error) {
	var comics Comics

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(Comics), nil
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

func SaveComicsCache(path string, comics Comics) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	CachePath := path + "Cache"
	err = os.WriteFile(CachePath, data, 0644)
	if err != nil {
		return err
	}
	err = os.Rename(CachePath, path)
	if err != nil {
		return err
	}

	return err
}
