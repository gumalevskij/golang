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

func GetLastComic(baseURL string) (int, error) {
	url := fmt.Sprintf("%s/info.0.json", baseURL)
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return -1, err
	}

	return comic.Num, nil
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

func FetchComics(baseURL string, existedComicsNumber int, limit int) ([]Comic, error) {
	lastComicNumber, err := GetLastComic(baseURL)
	if err != nil {
		return nil, err
	}

	if limit < 0 || limit > lastComicNumber {
		limit = lastComicNumber
	}

	var comics = make([]Comic, limit-existedComicsNumber)
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

		if i%100 == 0 {
			//fmt.Printf("\rFetched %d comics so far...", i-existedComicsNumber)
		}
	}

	return comics, nil
}
