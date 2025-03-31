ROOT := ./
BIN  := $(ROOT)bin/
SRC  := $(ROOT)src/

.PHONY: all clean

all: scraper xlsxwriter

scraper:
	go build -o $(BIN)scraper $(SRC)scraper/main.go

xlsxwriter:
	go build -o $(BIN)xlsxwriter $(SRC)xlsxwriter/main.go

clean:
	rm -rf $(BIN)
