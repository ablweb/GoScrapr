# GoScrapr
Scraper written in go
Convert html table into excel

# how to use
1. compile
```bash
make
```
2. run with
Scraper
```bash
./bin/scraper <url> [rule set]
```
The rule set is a json file, look at ruleSet.json for the structure 

wikixls
```bash
./bin/wikixls <url>
```
Convert html to excel

## TODO
0. rename wikixlsx, to work with general html table
1. Add test for scraping
2. write wikixlsx, for simple html wiki table to xlsx table
3. write same for html table, might be the same idk
