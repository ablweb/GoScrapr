package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)


func main() {
	
	scrapUrl := "https://en.wikipedia.org/wiki/Go_(programming_language)"
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	// go to wikipedia's golang page copy the first paragraphe and paste it in the terminal
	count := 0
	c.OnHTML("p", func(e *colly.HTMLElement) {
		if 0 < count && count < 2{
			fmt.Println(e.Text)
		}
		count++
	})


	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error while scraping %s\n", err.Error())
	})

	c.Visit(scrapUrl)
}
