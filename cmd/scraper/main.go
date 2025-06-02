package main

import (
	"fmt"
	"io"
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

func run(out, errOut io.Writer, args []string) int {
	// No URL
	if len(args) < 1 {
		fmt.Fprintln(errOut, "Error: No URL to scrap")
		fmt.Fprint(out, HELP)
		return 1
	}
	url := args[0]

	// Not valid URL
	reach, err := scraper.IsReachable(url)
	if !reach {
		fmt.Fprintln(errOut, "Error: URL is not reachable : ", err)
		return 1
	}

	var rules scraper.RuleSet
	if len(args) > 1 {
		rulesPath := args[1]
		rules, err = scraper.LoadRules(rulesPath)
		if err != nil {
			fmt.Fprintln(errOut, "Error: While loading rules", err)
			return 1
		}
	}

	scraped, err := scraper.ScrapMatches(url, rules)
	if err != nil {
		fmt.Fprintln(errOut, "Error: While scrapping", err)
		return 1
	}
	for _, e := range scraped {
		fmt.Fprintf(out, "'%s | %s'\n", e.Path, e.Content)
	}
	return 0
}

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args[1:]))
}
