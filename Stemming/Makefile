all: build

BINARY_NAME=myapp

build: clean
	go build -o $(BINARY_NAME)

run:
	./$(BINARY_NAME) -s "i'll follow you as long as you are following me"

run2:
	./$(BINARY_NAME) -s "follower brings bunch of questions"

run3:
	./$(BINARY_NAME) -s "Hannah's kitten"

clean:
	rm  -f $(BINARY_NAME)
