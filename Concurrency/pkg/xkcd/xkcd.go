package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Comic struct {
	Num            int
	Img            string
	Transcript     string
	Alt            string
	NormalizedText []string
}

func GetLastComicBinary(baseURL string, low int, high int) (int, error) {
	var lastComic Comic

	for low <= high {
		mid := low + (high-low)/2
		url := fmt.Sprintf("%s/%d/info.0.json", baseURL, mid)
		resp, err := http.Get(url)
		if err != nil {
			high = mid - 1
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if err := json.NewDecoder(resp.Body).Decode(&lastComic); err != nil {
				return -1, err
			}
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	if lastComic.Num != 0 {
		return lastComic.Num, nil
	}

	return -1, fmt.Errorf("unable to find the last comic")
}

func FetchComics(baseURL string, comicsNumber int) (*Comic, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", baseURL, comicsNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return nil, err
	}

	return &comic, nil
}
