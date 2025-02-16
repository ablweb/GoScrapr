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

var Reset = "\033[0m"
var Error = "\033[31m"
var Succes = "\033[32m"
var Warning = "\033[33m"
var Info = "\033[34m"

type Rule struct {
	Query    string `json:"query"`
	Priority int    `json:"priority"`
}
type RuleSet []Rule

func main() {
	// No URL
	if len(os.Args) < 2 {
		logError("No URL to scrap")
		fmt.Fprint(os.Stderr, HELP)
		os.Exit(1)
	}
	url := os.Args[1]
	// Not valid URL
	if !isReachable(&url) {
		fmt.Println(Error + "URL is not reachable" + Reset)
		os.Exit(1)
	}
	fmt.Println(Succes + "URL is reachable" + Reset)
	var rules RuleSet
	if len(os.Args) > 2 {
		rulesPath := os.Args[2]
		jsonContent, err_R := os.ReadFile(rulesPath)
		if err_R != nil {
			fmt.Printf(Error+"Error reading file: %s\n"+Reset, err_R)
			os.Exit(1)
		}
		err_U := json.Unmarshal(jsonContent, &rules)
		if err_U != nil {
			fmt.Printf(Error+"Error unmarshalling JSON: %s\n"+Reset, err_U)
			os.Exit(1)
		}
	} else {
		rules = nil
		fmt.Println(Warning + "No rules provided. Scraping URL without any rules." + Reset)
	}
	scraped, _ := scrap(&url, &rules)
	fmt.Print(scraped)
}

func logError(err string) {
	fmt.Fprintln(os.Stderr, "Error:", err)
}

func isReachable(url *string) bool {
	// check if url is valid
	fmt.Print(Info + "Checking URL: " + Reset)
	resp, err := http.Get(*url)
	if err != nil {
		print(Error + err.Error() + "\n" + Reset)
		return false
	} else {
		print(Info + fmt.Sprint(resp.StatusCode) + resp.Status + "\n" + Reset)
		return true
	}
}

func scrap(url *string, rules *RuleSet) (string, error) {
	c := colly.NewCollector()
	var dataMap = make(map[int]string)

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf(Info+"Visiting: %s\n"+Reset, r.URL)
	})
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
		fmt.Printf(Error+"Error while scraping: %s\n"+Reset, err.Error())
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
