package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type searchResults struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0",
}

func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func buildBingURLs(searchTerm string, pages int, count int) ([]string, error) {
	toScrape := []string{}

	query := url.QueryEscape(strings.TrimSpace(searchTerm))
	baseURL := "https://www.bing.com/search?q="

	for i := 0; i < pages; i++ {
		start := i * count
		scrapeURL := fmt.Sprintf("%s%s&first=%d", baseURL, query, start)
		toScrape = append(toScrape, scrapeURL)
	}

	return toScrape, nil
}

func bingResultParsing(response *http.Response, rank int) ([]searchResults, error) {
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	results := []searchResults{}
	rank++

	sel := doc.Find("li.b_algo")
	for i := range sel.Nodes {
		item := sel.Eq(i)
		link, _ := item.Find("a").Attr("href")
		title := item.Find("h3 a").Text()

		if link != "" && !strings.HasPrefix(link, "/") {
			results = append(results, searchResults{
				ResultRank:  rank,
				ResultURL:   link,
				ResultTitle: title,
			})
			rank++
		}
	}
	// html, _ := doc.Html()
	// fmt.Println(html[:500])

	return results, nil
}

func bingScrape(searchTerm string, proxyString interface{}, pages int, count int, backoff int) ([]searchResults, error) {
	results := []searchResults{}
	resultCounter := 0
	bingPages, err := buildBingURLs(searchTerm, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range bingPages {
		res, err := scrapeClientRequests(page, proxyString)
		if err != nil {
			return nil, err
		}
		data, err := bingResultParsing(res, resultCounter)
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

func getScrapeClient(proxyString interface{}) *http.Client {

	switch v := proxyString.(type) {

	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	default:
		return &http.Client{}
	}
}

func scrapeClientRequests(searchURL string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status: %d", res.StatusCode)
	}

	return res, nil
}

func main() {

	text := "Akshat Kumar"
	res, err := bingScrape(text, nil, 1, 30, 10)
	fmt.Println("Results count:", len(res))

	if err == nil {
		for _, res := range res {
			fmt.Println(res)
		}
	}
}