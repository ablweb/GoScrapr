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
	Query    string      `json:"query"`
	SubRules []QueryRule `json:"subrules,omitempty"`
	Priority int         `json:"priority"`
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
	Path    []string
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

	orderByPriority(&rules)
	applyRules(c, &rules, &matches)

	if err := c.Visit(url); err != nil {
		return nil, err
	}
	return matches, nil
}

// order rules in priotity order
func orderByPriority(rules *RuleSet) {
	sort.Slice(*rules, func(i, j int) bool {
		return (*rules)[i].Priority < (*rules)[j].Priority
	})
}

func applyRules(c *colly.Collector, rules *RuleSet, matches *MatchSet) {
	// Fallback behaviour, match everything
	if len(*rules) == 0 {
		c.OnHTML("html", func(e *colly.HTMLElement) {
			*matches = append(*matches, QueryMatch{
				Path:    []string{"html"},
				Content: e.Text + "\n"})
		})
		return
	}

	for _, rule := range *rules {
		var scopeId int = 0
		c.OnHTML(rule.Query, func(scope *colly.HTMLElement) {
			var path []string
			path = append(path, fmt.Sprintf("%s_%d", scope.Name, scopeId))
			findMatches(scope.DOM, &rule, path, matches)
			scopeId++
		})
	}
}

func findMatches(sel *goquery.Selection, subRule *QueryRule, path []string, matches *MatchSet) {
	// If no children, this is a terminal rule — collect data
	if len(subRule.SubRules) == 0 {
		text := strings.TrimSpace(sel.Text())
		tag := goquery.NodeName(sel)
		*matches = append(*matches, QueryMatch{
			Path:    append(path, tag),
			Content: text,
		})
		return
	}

	// Otherwise, recurse into each child
	for _, rule := range subRule.SubRules {
		sel.Find(rule.Query).Each(func(index int, s *goquery.Selection) {
			newPath := append(path, fmt.Sprintf("%s_%d", rule.Query, index))
			findMatches(s, &rule, newPath, matches)
		})
	}
}
