package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"time"
)

var googleDomains = map[string]string{

}
type searchResults struct{

	ResultRank int
	ResultURL string
	ResultTitle string
	ResultDesc string
}

var userAgents = []string{

}

func randomUserAgent() string{
	randNum := rand.Int() % len(userAgents)

	return userAgents[randNum]
}

func buildGoogleURLs()([]string,error){
	toScrape := []string{}

}

func googleScrape(string)([]searchResults,error){
	results:= []searchResults{}
	resultCounter := 0
	googlePages,err := buildGoogleURLs()

}

func main() {
	text := "Akshat Kumar"
	res,err := googleScrape(text)

	if err == nil{
		for _, res:= range res{
			fmt.Println(res)
		}
	}
}
