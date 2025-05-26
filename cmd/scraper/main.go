package main

import (
	"fmt"
	"os"

	"github.com/ablweb/GoScrapr/pkg/scraper"
)

const HELP = `Usage: scraper [URL] [RuleSet Path (optional)]

Parameters:
  URL            The website URL to scrape (required).
  RuleSet Path   Path to a JSON file defining scraping rules (optional).

RuleSet Format (JSON array):
[
  { "query": "h1", "priority": 1 },
  { "query": "h2", "priority": 2 },
  { "query": "p.desc", "priority": 1 }
]
- "query": CSS selector for elements to scrape.
- "priority": Determines the order of output (lower values appear first).
`

func main() {
	// No URL
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No URL to scrap")
		fmt.Fprint(os.Stderr, HELP)
		os.Exit(1)
	}
	url := os.Args[1]

	// Not valid URL
	reach, err := scraper.IsReachable(url)
	if !reach {
		fmt.Fprintln(os.Stderr, "Error: URL is not reachable : ", err)
		os.Exit(1)
	}

	var rules scraper.RuleSet
	if len(os.Args) > 2 {
		rulesPath := os.Args[2]
		rules, err = scraper.LoadRules(rulesPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: While loading rules", err)
			os.Exit(1)
		}
	}

	scraped, err := scraper.Scrap(url, rules)
	if err != nil {
			fmt.Fprintln(os.Stderr, "Error: While scrapping", err)
			os.Exit(1)
	}
	for _, e := range scraped { 
		fmt.Fprint(os.Stdout, "'")
		fmt.Fprint(os.Stdout, e.Path)
		fmt.Fprint(os.Stdout, " | ")
		fmt.Fprint(os.Stdout, e.Content)
		fmt.Fprint(os.Stdout, "'")
		fmt.Println()
	}
}
