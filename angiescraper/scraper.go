package angiescraper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type scraper struct {
	scrapeUrl string
}

type AngieScraper interface {
	Scrape() (AngieStatus, error)
}

func NewAngieScraper(scrapeUrl string) AngieScraper {
	s := new(scraper)
	s.scrapeUrl = scrapeUrl
	return s
}

func (s *scraper) Scrape() (AngieStatus, error) {

	scrapeResult := AngieStatus{
		Up: false,
	}

	resp, err := http.DefaultClient.Get(s.scrapeUrl)
	if err != nil {
		return scrapeResult, err
	}

	if resp.StatusCode != 200 {
		return scrapeResult, fmt.Errorf("Angie return %d for GET %s", resp.StatusCode, s.scrapeUrl)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&scrapeResult)

	if err != nil {
		return scrapeResult, fmt.Errorf("Angie return %d for GET %s", resp.StatusCode, s.scrapeUrl)
	}
	scrapeResult.Up = true

	return scrapeResult, nil
}
