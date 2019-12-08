TARGET = solitude

BINDIR = bin

build:
	go build -o $(BINDIR)/$(TARGET)

run: build
	./$(BINDIR)/$(TARGET)

clean:
	rm -rf $(BINDIR)

.PHONY: build run clean
