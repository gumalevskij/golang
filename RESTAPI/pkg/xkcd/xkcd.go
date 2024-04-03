package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"xkcd-fetcher/pkg/database"
)

type Comic struct {
	Num        int
	Img        string
	Transcript string
}

func LoadConfig(path string) (*database.Config, error) {
	return database.LoadConfig(path)
}

func FetchComics(baseURL string, limit int) ([]Comic, error) {
	var comics []Comic
	for i := 1; i <= limit; i++ {
		url := fmt.Sprintf("%s/%d/info.0.json", baseURL, i)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var comic Comic
		if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
			return nil, err
		}

		comics = append(comics, comic)
	}

	return comics, nil
}