ROOT := ./
BIN  := $(ROOT)bin/
SRC  := $(ROOT)cmd/

.PHONY: all clean

all: scraper wikixls

scraper:
	go build -o $(BIN)scraper $(SRC)scraper/main.go

wikixls:
	go build -o $(BIN)wikixls $(SRC)wikixls/main.go

clean:
	rm -rf $(BIN)
