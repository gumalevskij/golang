all: build

BINARY_NAME=xkcd

build:
	go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

run: build
	./$(BINARY_NAME) -o -n 10

run1: build
	./$(BINARY_NAME) -n 10

clean:
	rm $(BINARY_NAME)
