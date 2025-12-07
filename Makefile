.PHONY: build run clean


BINARY=bin/api


build:
	go build -o $(BINARY) ./cmd/api


run: build
	./$(BINARY)


clean:
	rm -f $(BINARY)


test:
	go test ./... -v