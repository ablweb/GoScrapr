# GoScrapr
GoScrapr is a web scraper tool written in go

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
1. this is the core tool, there should be an api to use it, so no all prints should be replaced by error code except for the text scraped, which should be passed to the cout
2. rule priority functionnality
3. write xlsx writer, for simple html table to xlsx table
4. python ML to image recognition, then image to xlsx table
5. [[MAYBE]] model to recognise patterns and predict tables
