package database

import (
	"encoding/json"
	"os"
)

func SaveComics(path string, comics map[string]interface{}) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
