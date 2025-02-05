package main

import (
	"fmt"
	"os"
	"net/http"
	"github.com/gocolly/colly/v2"
)

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No URL to scrape")
		os.Exit(1)
	}
	url := os.Args[1]
	fmt.Println(url)
	if !isReachable(url) {
		fmt.Println("URL is not reachable")
		os.Exit(1)
	}
	fmt.Println("URL is reachable")
	scrapUrl := "https://en.wikipedia.org/wiki/Go_(programming_language)"
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		//fmt.Printf("Visiting %s\n", r.URL)
	})

	// go to wikipedia's golang page copy the first paragraphe and paste it in the terminal
	count := 0
	c.OnHTML("p", func(e *colly.HTMLElement) {
		if 0 < count && count < 2{
			//fmt.Println(e.Text)
		}
		count++
	})

	c.OnError(func(r *colly.Response, err error) {
		//fmt.Printf("Error while scraping %s\n", err.Error())
	})

	c.Visit(scrapUrl)
}
