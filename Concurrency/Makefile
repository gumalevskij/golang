all: build

BINARY_NAME=xkcd

build:
	go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

run: build
	./$(BINARY_NAME)

clean:
	rm $(BINARY_NAME)
