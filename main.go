package main

import (
	"fmt"
	"os"
	"net/http"
	"google.golang.org/protobuf/proto"
	"github.com/gocolly/colly/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No URL to scrape")
		os.Exit(1)
	}
	url := os.Args[1]
	if !isReachable(url) {
		fmt.Println("URL is not reachable")
		os.Exit(1)
	}
	fmt.Println("URL is reachable")
	var ruleSet map[string]int
	if len(os.Args) > 2 {
		rulePath := os.Args[2]
		fmt.Println(rulePath)
		ruleSet = parseRuleSet(rulePath)
	}
	scrap(url, ruleSet)
}

func parseRuleSet(string rulePath) map[string]int {

}

func isReachable(url string) bool {
	// check if url is valid
	resp, err := http.Get(url)
	if err != nil {
		print(err.Error()+"\n")
		return false
	} else {
		print(fmt.Sprint(resp.StatusCode)+resp.Status+"\n")
		return true
	}
}

func scrap(url string, ruleSet map[string]int) bool {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error while scraping %s\n", err.Error())
	})

	c.Visit(url)
	return true
}
