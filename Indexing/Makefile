all: build

BINARY_NAME=xkcd

build:
	go build -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

run: build
	./$(BINARY_NAME) -i -s "i'll follow you as long as you are following me"

clean:
	rm $(BINARY_NAME)
