TARGET = solitude

BINDIR = bin

build:
	@cd cmd/solitude && go build -o ../../$(BINDIR)/$(TARGET)

run: build
	@./${BINDIR}/${TARGET}

test: build
	@./test.sh
	@./test_err.sh

clean:
	rm -rf $(BINDIR) tmp tmp.sl

.PHONY: build run test clean
