package main

import "github.com/gocolly/colly/v2"

func main() {
	// go to wikipedia's golang page copy the 10 first lines and paste it in the terminal
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
}
