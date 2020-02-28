TARGET = solitude
BINDIR = bin

NAME=aratanvm/solitude
VERSION=0.0.1

build:
	@cd cmd/solitude && go build -o ../../$(BINDIR)/$(TARGET)

run: build
	@./${BINDIR}/${TARGET}

test: build
	@./test.sh
	@./test_err.sh

clean:
	rm -rf $(BINDIR) tmp tmp.sl

docker/build:
	docker build -t $(NAME):$(VERSION) .

docker/run: docker/build
	docker run -it --name "Solitude" $(NAME):$(VERSION) /bin/bash

.PHONY: build run test clean
