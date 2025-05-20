package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type QueryRule struct {
	Scope    string `json:"scope"`
	Query    string `json:"query"`
	Priority int    `json:"priority"`
}
type RuleSet []QueryRule

// checks if a URL is valid and online
func IsReachable(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		err := fmt.Errorf("HTTP error %d", resp.StatusCode)
		return false, err
	}
	return true, nil
}

// load json rules
func LoadRules(rulesPath string) (RuleSet, error) {
	data, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, err
	}
	var rules RuleSet
	err = json.Unmarshal(data, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

type QueryMatch struct {
	Scope   string
	Tag     string
	Content string
}
type MatchSet []QueryMatch

// scap url with rule set and return the matches
func Scrap(url string, rules RuleSet) (MatchSet, error) {
	reach, err := IsReachable(url)
	if !reach {
		return nil, err
	}
	c := colly.NewCollector()
	var matches MatchSet

	// order rules in priotity order
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	// Fallback behaviour, match everything
	if len(rules) == 0 {
		c.OnHTML("html", func(e *colly.HTMLElement) {
			matches = append(matches, QueryMatch{
				Scope:   "html",
				Tag:     "html",
				Content: e.Text + "\n"})
		})
	// Normal behavou, match rules in given scopes
	} else {
		for _, rule := range rules {
			var queryScope = rule.Scope
			if queryScope == "" {
				queryScope = "html"
			}
			var scopeId int = 0
			c.OnHTML(queryScope, func(s *colly.HTMLElement) {
				s.DOM.Find(rule.Query).Each(
					func(_ int, m *goquery.Selection) {
						mtext := strings.TrimSpace(m.Text())
						mtag := goquery.NodeName(m)

						matches = append(matches, QueryMatch{
							Scope: fmt.Sprintf("%s_%d", s.Name,scopeId),
							Tag: mtag,
							Content: mtext})
					})
				scopeId++
			})
		}
	}
	if err := c.Visit(url); err != nil {
		return nil, err
	}
	return matches, nil
}
