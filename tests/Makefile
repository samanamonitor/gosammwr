
SRC_FILES=$(wildcard *.go)
BIN_FILES=$(patsubst %.go,bin/%,$(SRC_FILES))

all: $(BIN_FILES)

bin:
	mkdir -p $@

bin/%: %.go bin
	go build -o $@ $<

test: test.go
	go build $< 

clean:
	rm -Rf bin
