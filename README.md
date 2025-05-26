# GoScrapr
Scraper written in go
Convert html table into excel

# how to use
1. compile
```bash
make
```
2. run with
```bash
./bin/scraper <url> [rule set]
```
The rule set is a json file, look at ruleSet.json for the structure 

## TODO
0. rename wikixlsx, to work with general html table
1. implement a tree structure for matches
1. scrap takes the content instead of the url
2. Add test for scraping
3. write wikixlsx, for simple html wiki table to xlsx table
4. write same for html table, might be the same idk
