package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var googleDomains = map[string]string{}

type searchResults struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{}

func randomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := r % len(userAgents)

	return userAgents[randNum]
}

func buildGoogleURLs(searchTerm string, languageCode string, countryCode string, pages int, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBase, searchTerm, count, languageCode, start)
			toScrape = append(toScrape, scrapeURL)
		}
	} else {
		err := fmt.Errorf("country code is currently not supported")
		return nil, err
	}
	return toScrape, nil
}

func googleScrape(searchTerm, languageCode string, countryCode string, proxyString interface{}, pages, count, backoff int) ([]searchResults, error) {
	results := []searchResults{}
	resultCounter := 0
	googlePages, err := buildGoogleURLs(searchTerm, languageCode, countryCode, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := scrapeClientRequests(page, proxyString)
		if err != nil {
			return nil, err
		}
		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)
		for _, res := range data {
			results = append(results, res)
		}
		time.Sleep(time.Duration(backoff) * time.Second)

	}
	return results, nil

}

func main() {
	text := "Akshat Kumar"
	cc := "com"
	res, err := googleScrape(text, "en", cc, nil, 1, 30, 10)

	if err == nil {
		for _, res := range res {
			fmt.Println(res)
		}
	}
}
