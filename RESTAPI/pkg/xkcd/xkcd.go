package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Comic struct {
	Num        int
	Img        string
	Transcript string
	Alt        string
}

func GetLastComic(baseURL string) (Comic, error) {
    url := fmt.Sprintf("%s/info.0.json", baseURL)
    resp, err := http.Get(url)
    if err != nil {
        return Comic{}, err
    }
    defer resp.Body.Close()

    var comic Comic
    if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
        return Comic{}, err
    }

    return comic, nil
}

func FetchComics(baseURL string, existedComicsNumber int, limit int) ([]Comic, error) {
	lastComic, err := GetLastComic(baseURL)
    if err != nil {
        return nil, err
    }

	if limit < 0 || limit > lastComic.Num {
        limit = lastComic.Num
    }

	var comics []Comic
	for i := existedComicsNumber + 1; i <= limit; i++ {
		url := fmt.Sprintf("%s/%d/info.0.json", baseURL, i)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		var comic Comic
		if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
			return nil, err
		}

		comics = append(comics, comic)
	}

	return comics, nil
}
