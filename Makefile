TARGET = solitude

BINDIR = bin

build:
	go build -o $(BINDIR)/$(TARGET)

run: build
	./run.sh

test: build
	./test.sh

clean:
	rm -rf $(BINDIR) tmp.ll

.PHONY: build run test clean
