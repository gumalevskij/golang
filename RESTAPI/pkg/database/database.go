package database

import (
	"encoding/json"
	"io/ioutil"
)

func SaveComics(path string, comics map[string]interface{}) error {
	data, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
