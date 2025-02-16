package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/gocolly/colly/v2"
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

type Rule struct {
	Query    string `json:"query"`
	Priority int    `json:"priority"`
}
type RuleSet []Rule

func main() {
	// No URL
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No URL to scrap")
		fmt.Fprint(os.Stderr, HELP)
		os.Exit(1)
	}
	url := os.Args[1]
	// Not valid URL
	if !isReachable(&url) {
		fmt.Fprintln(os.Stderr, "Error: URL is not reachable")
		os.Exit(1)
	}
	var rules RuleSet
	if len(os.Args) > 2 {
		rulesPath := os.Args[2]
		jsonContent, err_R := os.ReadFile(rulesPath)
		if err_R != nil {
			fmt.Fprintln(os.Stderr, "Error: while reading file: ", err_R)
			os.Exit(1)
		}
		err_U := json.Unmarshal(jsonContent, &rules)
		if err_U != nil {
			fmt.Fprintln(os.Stderr, "Error: unmarshalling JSON: ", err_U)
			os.Exit(1)
		}
	} else {
		rules = nil
	}
	scraped, _ := scrap(&url, &rules)
	fmt.Fprintln(os.Stdout, scraped)
}

func isReachable(url *string) bool {
	// check if url is valid
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	} else {
		if (resp.StatusCode == 404) {
			fmt.Fprintln(os.Stderr, resp.Status)
			return false
		} else {
			return true
		}
	}
}

func scrap(url *string, rules *RuleSet) (string, error) {
	c := colly.NewCollector()
	var dataMap = make(map[int]string)

	if len(*rules) == 0 {
		// Scrap everything
		c.OnHTML("html", func(e *colly.HTMLElement) {
			dataMap[0] += e.Text + "\n"
		})
	} else {
		// Scrap with rules
		for _, rule := range *rules {
			c.OnHTML(rule.Query, func(e *colly.HTMLElement) {
				dataMap[rule.Priority] += e.Text + "\n"
			})
		}
	}
	c.OnError(func(r *colly.Response, err error) {
		fmt.Fprintln(os.Stderr, "Error : while scraping: ", err.Error())
	})

	if err := c.Visit(*url); err != nil {
		return "", err
	}

	var scraped string
	var keys []int
	for k := range dataMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		scraped += dataMap[k]
	}
	return scraped, nil
}
